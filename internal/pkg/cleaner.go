package pkg

import (
	"os"
	"time"
)

// RemoveDirectory synchronously removes a folder using os.RemoveAll which is
// usually the fastest way from Go on modern filesystems.
// Returns the duration it took.
func RemoveDirectory(targetPath string) (time.Duration, error) {
	start := time.Now()

	err := os.RemoveAll(targetPath)
	if err != nil {
		return 0, err
	}

	return time.Since(start), nil
}
