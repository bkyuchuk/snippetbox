package main

import (
	"testing"
	"time"

	"snippetbox.bogdandev.de/internal/assert"
)

func TestToHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			"UTC",
			time.Date(2026, 07, 18, 12, 45, 0, 0, time.UTC),
			"18 Jul 2026 at 12:45",
		},
		{
			"Empty",
			time.Time{},
			"",
		},
		{
			"CET",
			time.Date(2026, 07, 18, 12, 45, 0, 0, time.FixedZone("CET", 1*60*60)),
			"18 Jul 2026 at 11:45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := toHumanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}
