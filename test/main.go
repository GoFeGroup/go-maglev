package main

import (
	"fmt"

	"github.com/GoFeGroup/maglev"
)

func main() {
	bk := map[string]*maglev.Backend{
		"192.168.1.100": {
			ID:     100000100,
			Weight: 15,
		},
		"192.168.1.200": {
			ID:     100000200,
			Weight: 10,
		},
		"192.168.1.300": {
			ID:     100000300,
			Weight: 10,
		},
	}

	tb := maglev.GetLookupTable(bk, 100)

	for _, v := range tb {
		fmt.Println(v)
	}
}
