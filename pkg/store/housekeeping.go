package store

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/alice-lg/alice-lg/pkg/config"
)

// StartHousekeeping is a background task flushing
// memory and expireing caches.
func StartHousekeeping(ctx context.Context, cfg *config.Config) {

	for {
		if cfg.Housekeeping.Interval > 0 {
			time.Sleep(time.Duration(cfg.Housekeeping.Interval) * time.Minute)
		} else {
			time.Sleep(5 * time.Minute)
		}

		log.Println("Housekeeping started")

		// Expire the caches
		log.Println("Expiring caches")
		for _, source := range cfg.Sources {
			count := source.getInstance().ExpireCaches()
			log.Println("Expired", count, "entries for source", source.Name)
		}

		if cfg.Housekeeping.ForceReleaseMemory {
			// Trigger a GC and SCVG run
			log.Println("Freeing memory")
			debug.FreeOSMemory()
		}

		// Check if our services are still required
		select {
		case <-ctx.Done():
			fmt.Println("shutting down Housekeeping...")
			return
		default:
		}
	}
}
