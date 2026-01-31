package utils

import (
	"os"
	"path/filepath"
	"time"

	"github.com/OkanUysal/go-logger"
)

// CleanupOldFiles removes files older than the specified duration
func CleanupOldFiles(dir string, maxAge time.Duration) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil // Directory doesn't exist, nothing to clean
	}

	now := time.Now()
	removed := 0

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == dir {
			return nil
		}

		// Check if file/folder is older than maxAge
		if now.Sub(info.ModTime()) > maxAge {
			if err := os.RemoveAll(path); err != nil {
				logger.Warn("Failed to remove old file", logger.String("path", path), logger.Err(err))
				return nil // Continue despite error
			}
			removed++
			logger.Debug("Removed old file", logger.String("path", path))

			// If we removed a directory, skip its contents
			if info.IsDir() {
				return filepath.SkipDir
			}
		}

		return nil
	})

	if removed > 0 {
		logger.Info("Cleanup completed", logger.Int("removed", removed), logger.String("dir", dir))
	}

	return err
}

// StartCleanupRoutine starts a background goroutine that periodically cleans up old files
func StartCleanupRoutine(dir string, maxAge, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Run cleanup immediately on start
		if err := CleanupOldFiles(dir, maxAge); err != nil {
			logger.Error("Cleanup failed", logger.String("dir", dir), logger.Err(err))
		}

		// Then run periodically
		for range ticker.C {
			if err := CleanupOldFiles(dir, maxAge); err != nil {
				logger.Error("Cleanup failed", logger.String("dir", dir), logger.Err(err))
			}
		}
	}()
}
