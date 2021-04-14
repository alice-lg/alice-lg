package main

import (
	"log"
	"runtime/debug"
	"time"
)

func Housekeeping(config *Config) {
	for {
		if config.Housekeeping.Interval > 0 {
			time.Sleep(time.Duration(config.Housekeeping.Interval) * time.Minute)
		} else {
			time.Sleep(5 * time.Minute)
		}

		log.Println("Housekeeping started")

		// Expire the caches
		log.Println("Expiring caches")
		for _, source := range config.Sources {
			count := source.getInstance().ExpireCaches()
			log.Println("Expired", count, "entries for source", source.Name)
		}

		if config.Housekeeping.ForceReleaseMemory {
			// Trigger a GC and SCVG run
			log.Println("Freeing memory")
			debug.FreeOSMemory()
		}

	}
}
