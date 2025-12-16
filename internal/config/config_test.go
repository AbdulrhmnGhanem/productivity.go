package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanDatabaseID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Raw ID",
			input:    "a0e3e448792a4aa59f0d4576333457e9",
			expected: "a0e3e448792a4aa59f0d4576333457e9",
		},
		{
			name:     "Full URL with Query Params",
			input:    "https://www.notion.so/username/a0e3e448792a4aa59f0d4576333457e9?v=f291b0e4b2f64b7d818fe996318ecdf1",
			expected: "a0e3e448792a4aa59f0d4576333457e9",
		},
		{
			name:     "Full URL with Name and ID",
			input:    "https://www.notion.so/My-Reading-List-a0e3e448792a4aa59f0d4576333457e9",
			expected: "a0e3e448792a4aa59f0d4576333457e9",
		},
		{
			name:     "ID with Query Params only",
			input:    "a0e3e448792a4aa59f0d4576333457e9?v=123",
			expected: "a0e3e448792a4aa59f0d4576333457e9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanDatabaseID(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
