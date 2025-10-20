package main

import (
	"testing"

	"github.com/chickenzord/portosync/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestParseKseiAccountsWithName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]server.Account
	}{
		{
			name:  "single account",
			input: "personal:user1@example.com:password1",
			expected: map[string]server.Account{
				"personal": {
					Username: "user1@example.com",
					Password: "password1",
				},
			},
		},
		{
			name:  "multiple accounts",
			input: "personal:user1@example.com:password1,business:user2@example.com:password2",
			expected: map[string]server.Account{
				"personal": {
					Username: "user1@example.com",
					Password: "password1",
				},
				"business": {
					Username: "user2@example.com",
					Password: "password2",
				},
			},
		},
		{
			name:  "accounts with spaces",
			input: " personal : user1@example.com : password1 , business : user2@example.com : password2 ",
			expected: map[string]server.Account{
				"personal": {
					Username: "user1@example.com",
					Password: "password1",
				},
				"business": {
					Username: "user2@example.com",
					Password: "password2",
				},
			},
		},
		{
			name:  "password with colon",
			input: "personal:user1@example.com:pass:word:123",
			expected: map[string]server.Account{
				"personal": {
					Username: "user1@example.com",
					Password: "pass:word:123",
				},
			},
		},
		{
			name:     "empty string",
			input:    "",
			expected: map[string]server.Account{},
		},
		{
			name:     "invalid format - only two parts",
			input:    "personal:user1@example.com",
			expected: map[string]server.Account{},
		},
		{
			name:     "invalid format - only one part",
			input:    "personal",
			expected: map[string]server.Account{},
		},
		{
			name:  "empty password",
			input: "personal:user1@example.com:",
			expected: map[string]server.Account{
				"personal": {
					Username: "user1@example.com",
					Password: "",
				},
			},
		},
		{
			name:  "empty username",
			input: "personal::password1",
			expected: map[string]server.Account{
				"personal": {
					Username: "",
					Password: "password1",
				},
			},
		},
		{
			name:     "empty name",
			input:    ":user1@example.com:password1",
			expected: map[string]server.Account{},
		},
		{
			name:  "special characters in password",
			input: "personal:user1@example.com:p@ss$w0rd!#%",
			expected: map[string]server.Account{
				"personal": {
					Username: "user1@example.com",
					Password: "p@ss$w0rd!#%",
				},
			},
		},
		{
			name:  "trailing comma",
			input: "personal:user1@example.com:password1,",
			expected: map[string]server.Account{
				"personal": {
					Username: "user1@example.com",
					Password: "password1",
				},
			},
		},
		{
			name:  "multiple commas",
			input: "personal:user1@example.com:password1,,business:user2@example.com:password2",
			expected: map[string]server.Account{
				"personal": {
					Username: "user1@example.com",
					Password: "password1",
				},
				"business": {
					Username: "user2@example.com",
					Password: "password2",
				},
			},
		},
		{
			name:  "duplicate names - last one wins",
			input: "personal:user1@example.com:password1,personal:user2@example.com:password2",
			expected: map[string]server.Account{
				"personal": {
					Username: "user2@example.com",
					Password: "password2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseKseiAccountsWithName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseKseiAccountsWithName_EmptyResult(t *testing.T) {
	result := parseKseiAccountsWithName("")
	assert.NotNil(t, result)
	assert.Empty(t, result)
}
