package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

// CPU性能分析
func main() {
	var isCPUPprof bool
	flag.BoolVar(&isCPUPprof, "cpu", false, "turn cpu pprof on")
	flag.Parse()

	if isCPUPprof {
		file, err := os.Create("./cpu.pprof")
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(file)
		defer pprof.StopCPUProfile()
	}

	joinSlice()
}

// 一段性能问题代码
func joinSlice() []int {
	num := 10000000
	//var slices []int
	slices := make([]int, 0, num)
	for i := 0; i < num; i++ {
		// 故意造成多次切片添加，由于每次操作可能会有内存重新分配和移动，性能较低
		slices = append(slices, i+1)
	}
	return slices
}
