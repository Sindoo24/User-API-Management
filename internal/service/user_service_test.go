package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name     string
		dob      time.Time
		expected int
	}{
		{
			name:     "Age 34 - birthday already passed this year",
			dob:      time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC),
			expected: 34,
		},
		{
			name:     "Age 33 - birthday not yet this year",
			dob:      time.Date(1991, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: 33,
		},
		{
			name:     "Age 0 - born this year",
			dob:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "Age 25 - leap year birth",
			dob:      time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			expected: 24,
		},
		{
			name:     "Age 100 - very old person",
			dob:      time.Date(1924, 6, 15, 0, 0, 0, 0, time.UTC),
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateAge(tt.dob)

			if result < tt.expected-1 || result > tt.expected+1 {
				t.Errorf("calculateAge(%v) = %d; want approximately %d", tt.dob, result, tt.expected)
			}
		})
	}
}

func TestCalculateAgeEdgeCases(t *testing.T) {
	now := time.Now()

	t.Run("Birthday today", func(t *testing.T) {
		dob := now.AddDate(-30, 0, 0)
		age := calculateAge(dob)
		if age != 30 {
			t.Errorf("calculateAge for birthday today = %d; want 30", age)
		}
	})

	t.Run("Birthday tomorrow", func(t *testing.T) {
		dob := now.AddDate(-30, 0, 1)
		age := calculateAge(dob)
		if age != 29 {
			t.Errorf("calculateAge for birthday tomorrow = %d; want 29", age)
		}
	})

	t.Run("Birthday yesterday", func(t *testing.T) {
		dob := now.AddDate(-30, 0, -1)
		age := calculateAge(dob)
		if age != 30 {
			t.Errorf("calculateAge for birthday yesterday = %d; want 30", age)
		}
	})
}
