//go:build !linux

package pkg

import "time"

func accessTime(_ string) (time.Time, bool) {
	return time.Time{}, false
}
