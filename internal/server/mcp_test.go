package server

import (
	"context"
	"testing"

	"github.com/chickenzord/goksei"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

func TestListAccountNamesResult_Description(t *testing.T) {
	tests := []struct {
		name     string
		result   ListAccountNamesResult
		expected string
	}{
		{
			name: "multiple accounts",
			result: ListAccountNamesResult{
				AccountNames: []string{"personal", "business", "family"},
			},
			expected: "Available accounts: personal, business, family",
		},
		{
			name: "single account",
			result: ListAccountNamesResult{
				AccountNames: []string{"personal"},
			},
			expected: "Available accounts: personal",
		},
		{
			name: "no accounts",
			result: ListAccountNamesResult{
				AccountNames: []string{},
			},
			expected: "No accounts configured",
		},
		{
			name: "nil accounts",
			result: ListAccountNamesResult{
				AccountNames: nil,
			},
			expected: "No accounts configured",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.result.Description()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMCP_handleListAccountNames(t *testing.T) {
	// Create a mock MCP server with test accounts
	accounts := map[string]Account{
		"personal": {
			Username: "user1@example.com",
			Password: "password1",
		},
		"business": {
			Username: "user2@example.com",
			Password: "password2",
		},
	}

	// Create mock KSEI clients
	kseiClients := make(map[string]*goksei.Client)
	for name := range accounts {
		kseiClients[name] = &goksei.Client{}
	}

	mcpServer := &MCP{
		kseiClients: kseiClients,
	}

	ctx := context.Background()
	req := &mcp.CallToolRequest{}
	args := ListAccountNamesArgs{}

	result, data, err := mcpServer.handleListAccountNames(ctx, req, args)

	// Verify no error occurred
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, data)

	// Verify the result contains both account names
	assert.Len(t, data.AccountNames, 2)
	assert.Contains(t, data.AccountNames, "personal")
	assert.Contains(t, data.AccountNames, "business")

	// Verify the result content
	assert.Len(t, result.Content, 1)
	textContent, ok := result.Content[0].(*mcp.TextContent)
	assert.True(t, ok)
	assert.Contains(t, textContent.Text, "Available accounts:")
	assert.Contains(t, textContent.Text, "personal")
	assert.Contains(t, textContent.Text, "business")
}

func TestMCP_handleListAccountNames_NoAccounts(t *testing.T) {
	mcpServer := &MCP{
		kseiClients: make(map[string]*goksei.Client),
	}

	ctx := context.Background()
	req := &mcp.CallToolRequest{}
	args := ListAccountNamesArgs{}

	result, data, err := mcpServer.handleListAccountNames(ctx, req, args)

	// Verify no error occurred
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, data)

	// Verify the result contains no account names
	assert.Len(t, data.AccountNames, 0)

	// Verify the result content shows "No accounts configured"
	assert.Len(t, result.Content, 1)
	textContent, ok := result.Content[0].(*mcp.TextContent)
	assert.True(t, ok)
	assert.Equal(t, "No accounts configured", textContent.Text)
}
