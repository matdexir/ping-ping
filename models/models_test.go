package models_test

import (
	// "strings"
	"testing"
	"time"

	"github.com/matdexir/ping-ping/models"
)

func TestSerialize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		items    []models.SerializableItem
		expected string
	}{
		{
			name:     "Serialize platforms",
			items:    []models.SerializableItem{models.ANDROID, models.IOS},
			expected: "android, iOS",
		},
		{
			name:     "Serialize countries",
			items:    []models.SerializableItem{models.Taiwan, models.USA},
			expected: "TW, US",
		},
		{
			name:     "Serialize genders",
			items:    []models.SerializableItem{models.MALE, models.FEMALE},
			expected: "M, F",
		},
		{
			name:     "Empty slice",
			items:    []models.SerializableItem{},
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := models.Serialize(tc.items)
			if actual != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestSponsoredPostValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		sp      *models.SponsoredPost
		wantErr bool
		errStr  string // Optional expected error message substring
	}{
		{
			name:    "valid post",
			sp:      &models.SponsoredPost{Title: "Test Post", StartAt: time.Now(), EndAt: time.Now().Add(time.Hour)},
			wantErr: false,
		},
		{
			name:    "empty title",
			sp:      &models.SponsoredPost{StartAt: time.Now(), EndAt: time.Now().Add(time.Hour)},
			wantErr: true,
			errStr:  "required",
		},
		{
			name:    "missing startAt",
			sp:      &models.SponsoredPost{Title: "Test Post", EndAt: time.Now().Add(time.Hour)},
			wantErr: true,
			errStr:  "required",
		},
		{
			name:    "missing endAt",
			sp:      &models.SponsoredPost{Title: "Test Post", StartAt: time.Now()},
			wantErr: true,
			errStr:  "required",
		},
		{
			name:    "endAt before startAt",
			sp:      &models.SponsoredPost{Title: "Test Post", StartAt: time.Now().Add(time.Hour), EndAt: time.Now()},
			wantErr: true,
			errStr:  "must be after",
		},
		{
			name:    "invalid age range in settings",
			sp:      &models.SponsoredPost{Title: "Test Post", StartAt: time.Now(), EndAt: time.Now().Add(time.Hour), Conditions: models.Settings{AgeStart: 100, AgeEnd: 80}},
			wantErr: true,
			errStr:  "ltecsfield",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sp.Validate()
			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
