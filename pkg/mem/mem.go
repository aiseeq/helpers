package mem

import (
	"fmt"
	"runtime"
)

func GetStats() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return fmt.Sprintf("Alloc = %vK, TotalAlloc = %vK, Sys = %vK, NumGC = %v",
		m.Alloc/1024,
		m.TotalAlloc/1024,
		m.Sys/1024,
		m.NumGC)
}
