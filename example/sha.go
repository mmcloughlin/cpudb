package main

import (
	"fmt"

	"github.com/mmcloughlin/cpudb"
)

func main() {
	for _, cpu := range cpudb.CPUs {
		leaf, found := cpu.CPUID(7, 0)
		if !found {
			continue
		}

		if (leaf.EBX>>29)&1 == 0 {
			continue
		}

		fmt.Printf("%s: %s\n", cpu.Alias, cpu.Type)
	}
}
