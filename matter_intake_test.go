package main

import (
	"testing"
	"time"
)

func TestFollowUpNeeded(t *testing.T) {
	now := time.Date(2026, 9, 9, 10, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		deadline time.Time
		want     bool
	}{
		{"inside window", now.Add(48 * time.Hour), true},
		{"outside window", now.Add(96 * time.Hour), false},
		{"past deadline", now.Add(-time.Hour), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := followUpNeeded(now, tc.deadline); got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}
