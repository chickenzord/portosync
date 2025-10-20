package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKseiAccountsString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:  "single account",
			input: "user1@example.com:password1",
			expected: map[string]string{
				"user1@example.com": "password1",
			},
		},
		{
			name:  "multiple accounts",
			input: "user1@example.com:password1,user2@example.com:password2",
			expected: map[string]string{
				"user1@example.com": "password1",
				"user2@example.com": "password2",
			},
		},
		{
			name:  "accounts with spaces",
			input: " user1@example.com : password1 , user2@example.com : password2 ",
			expected: map[string]string{
				"user1@example.com": "password1",
				"user2@example.com": "password2",
			},
		},
		{
			name:  "password with colon",
			input: "user1@example.com:pass:word:123",
			expected: map[string]string{
				"user1@example.com": "pass:word:123",
			},
		},
		{
			name:     "empty string",
			input:    "",
			expected: map[string]string{},
		},
		{
			name:     "invalid format - no colon",
			input:    "user1@example.com",
			expected: map[string]string{},
		},
		{
			name:  "invalid format - missing value",
			input: "user1@example.com:",
			expected: map[string]string{
				"user1@example.com": "",
			},
		},
		{
			name:  "mixed valid and invalid",
			input: "user1@example.com:password1,invalid,user2@example.com:password2",
			expected: map[string]string{
				"user1@example.com": "password1",
				"user2@example.com": "password2",
			},
		},
		{
			name:  "special characters in password",
			input: "user1@example.com:p@ss$w0rd!#%",
			expected: map[string]string{
				"user1@example.com": "p@ss$w0rd!#%",
			},
		},
		{
			name:  "trailing comma",
			input: "user1@example.com:password1,",
			expected: map[string]string{
				"user1@example.com": "password1",
			},
		},
		{
			name:  "multiple commas",
			input: "user1@example.com:password1,,user2@example.com:password2",
			expected: map[string]string{
				"user1@example.com": "password1",
				"user2@example.com": "password2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseKseiAccountsString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseKseiAccountsString_EmptyResult(t *testing.T) {
	result := parseKseiAccountsString("")
	assert.NotNil(t, result)
	assert.Empty(t, result)
}
