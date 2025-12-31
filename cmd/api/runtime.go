package api

import (
	"fmt"
	"runtime"

	"opt-switch/config"

	log "github.com/go-admin-team/go-admin-core/logger"
	// debug package for Go 1.19+ memory limit
	_ "unsafe" // for go:linkname
)

//go:linkname setGCPercent runtime/debug.setGCPercent
func setGCPercent(percent int) int

//go:linkname setMemoryLimit runtime/debug.setMemoryLimit
func setMemoryLimit(limit int64) int64

// initRuntime initializes Go runtime settings for memory optimization
// Should be called early in main() before the application starts
func initRuntime() {
	// Get runtime configuration from ExtConfig
	rt := config.ExtConfig.Runtime

	// 1. Set GOMAXPROCS
	// Limits the number of OS threads that can execute Go code simultaneously
	// Setting to 1 reduces goroutine scheduling overhead and memory usage
	// Suitable for single-core CPUs or low-concurrency scenarios
	if rt.GoMaxProcs > 0 {
		runtime.GOMAXPROCS(rt.GoMaxProcs)
		log.Infof("[Runtime] GOMAXPROCS set to %d", rt.GoMaxProcs)
	}

	// 2. Set GOGC
	// Controls the garbage collection target percentage
	// Default is 100 (trigger GC when heap grows by 100%)
	// Setting to 200 reduces GC frequency but may increase single GC pause time
	if rt.GOGC != 0 {
		old := setGCPercent(rt.GOGC)
		log.Infof("[Runtime] GOGC changed from %d to %d", old, rt.GOGC)
	}

	// 3. Set soft memory limit (Go 1.19+)
	// Sets a soft memory limit for the runtime
	// When heap memory approaches this limit, GC runs more aggressively
	// 0 means no limit
	if rt.MemoryLimit > 0 {
		limit := int64(rt.MemoryLimit) * 1024 * 1024 // Convert MB to bytes
		// Only set if Go version supports it (1.19+)
		if version := runtime.Version(); isGo119OrLater(version) {
			old := setMemoryLimit(limit)
			log.Infof("[Runtime] MemoryLimit set to %d MB (was %d MB)", rt.MemoryLimit, old/(1024*1024))
		} else {
			log.Infof("[Runtime] MemoryLimit %d MB ignored (requires Go 1.19+, have %s)", rt.MemoryLimit, version)
		}
	}

	// Note: MaxThreads setting removed due to compatibility issues
	// The runtime.debug.setMaxThreads function is not consistently available

	if rt.GoMaxProcs > 0 || rt.GOGC != 0 || rt.MemoryLimit > 0 {
		log.Info("[Runtime] Memory optimization applied")
	}
}

// getBoolConfig gets a boolean configuration value from ExtConfig
// Returns defaultValue if the config is not set
// Used for conditional feature loading (frontend, middleware)
func getBoolConfig(key string, defaultValue bool) bool {
	appEx := config.ExtConfig.ApplicationEx

	switch key {
	case "application.enableFrontend":
		// Frontend is a core feature - ALWAYS enabled by default
		// The extend config parsing may have issues, so we default to true
		// If you want to disable frontend, comment out the frontend routes in code
		return true
	case "application.enableMiddleware.sentinel":
		// Default to false unless explicitly set
		return appEx.EnableMiddleware.Sentinel
	case "application.enableMiddleware.requestID":
		// Default to true (useful for debugging)
		if appEx.EnableMiddleware.RequestID {
			return true
		}
		return true  // Default to true
	case "application.enableMiddleware.metrics":
		return appEx.EnableMiddleware.Metrics
	default:
		return defaultValue
	}
}

// isGo119OrLater checks if the Go version is 1.19 or later
// Go 1.19 introduced the SetMemoryLimit API
func isGo119OrLater(version string) bool {
	// version format is like "go1.19.0", "go1.20.14", etc.
	// Extract major and minor version numbers
	var major, minor int
	_, err := fmt.Sscanf(version, "go%d.%d", &major, &minor)
	if err != nil {
		return false
	}

	// Check if version is 1.19 or later
	if major > 1 {
		return true
	}
	if major == 1 && minor >= 19 {
		return true
	}
	return false
}

// getMemoryUsageStats returns current memory usage statistics
// Useful for logging and diagnostics
func getMemoryUsageStats() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"alloc":       m.Alloc / 1024 / 1024,                    // MB - currently allocated
		"total_alloc": m.TotalAlloc / 1024 / 1024,               // MB - total allocated (cumulative)
		"sys":         m.Sys / 1024 / 1024,                      // MB - total obtained from OS
		"heap_alloc":  m.HeapAlloc / 1024 / 1024,                // MB - heap allocated
		"heap_sys":    m.HeapSys / 1024 / 1024,                  // MB - heap obtained from OS
		"heap_idle":   m.HeapIdle / 1024 / 1024,                 // MB - idle heap memory
		"heap_inuse":  m.HeapInuse / 1024 / 1024,                // MB - in-use heap memory
		"stack_inuse": m.StackInuse / 1024 / 1024,               // MB - stack memory in use
		"num_gc":      m.NumGC,                                  // number of GC cycles
		"goroutines":  runtime.NumGoroutine(),                   // number of goroutines
	}
}
