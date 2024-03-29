package main

import (
	"FEEDFACE.COM/facade/gfx"
	"fmt"
	"runtime"
	"time"
)

const DEBUG_DIAG     = false

func MemUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	ret := ""

	ret += fmt.Sprintf("live:%d#  heap:%.1fmb  sys:%.1fmb  stack:%.1fmb",
		m.Mallocs-m.Frees, mib(m.HeapAlloc), mib(m.HeapSys), mib(m.StackInuse))

	//    ret += fmt.Sprintf("alloc %6.2f/%6.2f MiB   ",mib(m.Alloc),mib(m.TotalAlloc))
	//    ret += fmt.Sprintf("total %6.2f MiB  ",mib(m.TotalAlloc))
	//    ret += fmt.Sprintf("  sys %6.2f MiB  ",mib(m.Sys))
	//    ret += fmt.Sprintf("numgc %v",m.NumGC)
	return ret

}

func StartGC() {
	runtime.GC()
}

func mib(bits uint64) float64 { return float64(bits) / (1024. * 1024) }

var diagStart time.Time
var diagStats *gfx.RB

func DiagStart() {
	if !DEBUG_DIAG {
		return
	}
	diagStart = time.Now()
}

func DiagDone() {
	if !DEBUG_DIAG {
		return
	}
	delta := time.Since(diagStart)
	if diagStats == nil {
		diagStats = gfx.NewRB(10)
	}
	diagStats.Add(float32(delta.Seconds()))
}

func InfoDiag() string {
	if !DEBUG_DIAG {
		return ""
	}

	aspf := diagStats.Average()
	afps := 1. / aspf

	return fmt.Sprintf("avg %5.2fms (%.2ffps)", aspf*1000., afps)

}
