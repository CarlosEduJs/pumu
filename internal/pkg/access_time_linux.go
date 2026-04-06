//go:build linux

package pkg

import (
	"os"
	"syscall"
	"time"
)

func accessTime(path string) (time.Time, bool) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, false
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return time.Time{}, false
	}
	return time.Unix(stat.Atim.Sec, stat.Atim.Nsec), true
}
