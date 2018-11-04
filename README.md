# cpudb

[CPUID](https://en.wikipedia.org/wiki/CPUID) database derived from [InstLatx64](http://instlatx64.atw.hu).

[![GoDoc Reference](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](http://godoc.org/github.com/mmcloughlin/cpudb)
[![Build status](https://img.shields.io/travis/mmcloughlin/cpudb.svg?style=flat-square)](https://travis-ci.org/mmcloughlin/cpudb)

## Description

This package contains a list of CPUs from the [InstLatx64 database](https://github.com/InstLatx64/InstLatx64) along with their CPUID values. The [`CPU.CPUID()`](https://godoc.org/github.com/mmcloughlin/cpudb#CPU.CPUID) method enables you to perform CPUID queries on each processor, therefore enabling filtering to processors with specific capabilities.

## Example

The following [example](example/sha.go) shows how to search the database for CPUs supporting [SHA extensions](https://en.wikipedia.org/wiki/Intel_SHA_extensions).

[embedmd]:# (example/sha.go)
```go
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
```

Output:

[embedmd]:# (example/sha.out)
```out
Summit Ridge: OctalCore AMD Summit Ridge, 13174 MHz (132 x 100)
Summit Ridge: OctalCore AMD Ryzen 7 1700X, 3400 MHz (34 x 100)
Summit Ridge: QuadCore AMD Ryzen 5 1500X, 3600 MHz (36 x 100)
Summit Ridge: OctalCore AMD Ryzen 7 1800X, 3600 MHz (36 x 100)
Threadripper 2: 32-Core AMD Ryzen Threadripper 2990WX, 3000 MHz (30 x 100)
Pinnacle Ridge: OctalCore AMD Ryzen 7 2700X, 3700 MHz (37 x 100)
Raven Ridge: Mobile QuadCore AMD Ryzen 5 2500U, 2400 MHz (24 x 100)
Apollo Lake-D: QuadCore Intel Celeron J3455, 2200 MHz (22 x 100)
Apollo Lake: QuadCore , 466 MHz (24 x 19)
Denverton: 16-Core Intel Atom C3958, 2000 MHz (20 x 100)
Cannon Lake-U: DualCore Intel Core i3-8121U, 3200 MHz (32 x 100)
Gemini Lake-D: QuadCore Intel Celeron J4105, 2400 MHz (24 x 100)
```
