package main

import (
	"fmt"

	"github.com/GoFeGroup/maglev"
	"github.com/GoFeGroup/maglev/loadbalancer"
)

func main() {
	bk := map[string]*loadbalancer.Backend{
		"192.168.1.100": {
			FEPortName: "",
			ID:         100000100,
			Weight:     15,
			NodeName:   "100000100",
			// L3n4Addr:   loadbalancer.L3n4Addr{},
			State:     0,
			Preferred: false,
		},
		"192.168.1.200": {
			FEPortName: "",
			ID:         100000200,
			Weight:     10,
			NodeName:   "100000200",
			// L3n4Addr:   loadbalancer.L3n4Addr{},
			State:     0,
			Preferred: false,
		},
		"192.168.1.300": {
			FEPortName: "",
			ID:         100000300,
			Weight:     10,
			NodeName:   "100000300",
			// L3n4Addr:   loadbalancer.L3n4Addr{},
			State:     0,
			Preferred: false,
		},
	}

	tb := maglev.GetLookupTable(bk, 100)

	for _, v := range tb {
		fmt.Println(v)
	}
}
