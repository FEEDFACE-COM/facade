
package render

import (
    "runtime"
    "fmt"
)


func MemUsage() string {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    ret := ""
    ret += fmt.Sprintf("alloc %6.2f/%6.2f MiB   ",mib(m.Alloc),mib(m.TotalAlloc))
//    ret += fmt.Sprintf("total %6.2f MiB  ",mib(m.TotalAlloc))
    ret += fmt.Sprintf("  sys %6.2f MiB  ",mib(m.Sys))
    ret += fmt.Sprintf("numgc %v",m.NumGC)
    return ret    
    
}


func mib(bits uint64) float64 { return float64(bits) / (1024.*1024) }
