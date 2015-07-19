// Package tgflake is a wrapper around sonyflake
package tgflake

import (
	"time"

	"github.com/sony/sonyflake"
)

var (
	flakes = map[int64]map[string]*sonyflake.Sonyflake{}
)

// Flake returns the specific flake for a certain application
func Flake(appID int64, flake string) *sonyflake.Sonyflake {
	if f, ok := flakes[appID][flake]; ok {
		return f
	}

	if _, ok := flakes[appID]; !ok {
		flakes[appID] = map[string]*sonyflake.Sonyflake{}
	}

	var st sonyflake.Settings
	st.StartTime = time.Date(2014, 12, 17, 18, 7, 0, 0, time.UTC)
	flakes[appID][flake] = sonyflake.NewSonyflake(st)
	if flakes[appID][flake] == nil {
		panic("sonyflake not created")
	}

	return flakes[appID][flake]
}

// FlakeNextID will get the next ID from the
func FlakeNextID(appID int64, flake string) (uint64, error) {
	return Flake(appID, flake).NextID()
}

// RemoveAllFlakes removes all active flakes! Do not use this in production!
func RemoveAllFlakes() {
	flakes = map[int64]map[string]*sonyflake.Sonyflake{}
}
