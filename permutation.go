package maglev

import (
	"context"
	"encoding/base64"

	"github.com/GoFeGroup/maglev/murmur3"
	"github.com/GoFeGroup/maglev/workerpool"
)

var (
	seedMurmur      uint32
	DefaultHashSeed = "JLfvgnHc2kaSUFaI"
	permutation     []uint64
)

func init() {
	d, _ := base64.StdEncoding.DecodeString(DefaultHashSeed)
	seedMurmur = uint32(d[0])<<24 | uint32(d[1])<<16 | uint32(d[2])<<8 | uint32(d[3])
}

func getOffsetAndSkip(backend string, m uint64) (uint64, uint64) {
	h1, h2 := murmur3.Hash128([]byte(backend), seedMurmur)
	offset := h1 % m
	skip := (h2 % (m - 1)) + 1

	return offset, skip
}

func getPermutation(backends []string, m uint64, numCPU int) []uint64 {
	// The idea is to split the calculation into batches so that they can be
	// concurrently executed. We limit the number of concurrent goroutines to
	// the number of available CPU cores. This is because the calculation does
	// not block and is completely CPU-bound. Therefore, adding more goroutines
	// would result into an overhead (allocation of stackframes, stress on
	// scheduling, etc) instead of a performance gain.

	bCount := len(backends)
	if size := uint64(bCount) * m; size > uint64(len(permutation)) {
		// Reallocate slice so we don't have to allocate again on the next
		// call.
		permutation = make([]uint64, size)
	}

	batchSize := bCount / numCPU
	if batchSize == 0 {
		batchSize = bCount
	}

	// Since no other goroutine is controlling the WorkerPool, it is safe to
	// ignore the returned error from wp methods. Also as our task func never
	// return any error, we have no use returned value from Drain() and don't
	// need to provide an id to Submit().
	wp := workerpool.New(numCPU)
	defer wp.Close()
	for g := 0; g < bCount; g += batchSize {
		from, to := g, g+batchSize
		if to > bCount {
			to = bCount
		}
		wp.Submit("", func(_ context.Context) error {
			for i := from; i < to; i++ {
				offset, skip := getOffsetAndSkip(backends[i], m)
				permutation[i*int(m)] = offset % m
				for j := uint64(1); j < m; j++ {
					permutation[i*int(m)+int(j)] = (permutation[i*int(m)+int(j-1)] + skip) % m
				}
			}
			return nil
		})
	}
	wp.Drain()

	return permutation[:bCount*int(m)]
}
