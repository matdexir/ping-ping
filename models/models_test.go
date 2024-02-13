package models_test

import (
	"testing"

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
		{
			name:     "Mixed types",
			items:    []models.SerializableItem{models.ANDROID, models.Taiwan},
			expected: "android, TW", // Ignore invalid item
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
