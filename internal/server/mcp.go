package server

import (
	"context"
	"maps"
	"net/http"
	"strings"
	"time"

	"github.com/chickenzord/goksei"
	"github.com/chickenzord/portosync/internal/version"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCP wraps the MCP SDK server
type MCP struct {
	kseiClients map[string]*goksei.Client
	mcpServer   *mcp.Server
}

// selectKseiClients get clients by multiple names,
// if empty or nil, it will return all clients
func (m *MCP) selectKseiClients(names []string) map[string]*goksei.Client {
	if len(names) == 0 {
		clients := make(map[string]*goksei.Client)
		maps.Copy(clients, m.kseiClients)

		return clients
	}

	// Return clients matching the provided names
	clients := make(map[string]*goksei.Client)

	for _, name := range names {
		if client, ok := m.kseiClients[name]; ok {
			clients[name] = client
		}
	}

	return clients
}

func (m *MCP) getKseiClientNames() []string {
	names := make([]string, 0, len(m.kseiClients))
	for name := range m.kseiClients {
		names = append(names, name)
	}

	return names
}

func (m *MCP) RunHTTP(ctx context.Context, bindAddress string) error {
	httpHandler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return m.mcpServer
	}, nil)

	return http.ListenAndServe(bindAddress, httpHandler)
}

func (m *MCP) RunStdio(ctx context.Context) error {
	return m.mcpServer.Run(ctx, &mcp.StdioTransport{})
}

// NewMCP creates a new MCP server using the official MCP Go SDK
func NewMCP(accounts map[string]Account, plainPassword bool, authCacheDir string) *MCP {
	authStore, err := goksei.NewFileAuthStore(authCacheDir)
	if err != nil {
		panic(err)
	}

	gokseiClients := make(map[string]*goksei.Client, len(accounts))
	for name, account := range accounts {
		gokseiClients[name] = goksei.NewClient(goksei.ClientOpts{
			Username:      account.Username,
			Password:      account.Password,
			PlainPassword: plainPassword,
			Timeout:       1 * time.Minute,
			AuthStore:     authStore,
		})
	}

	s := &MCP{
		kseiClients: gokseiClients,
	}

	// Create MCP server with implementation info
	versionInfo := version.Get()
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "portosync",
		Version: versionInfo.Version,
		Title:   "PortoSync - Financial Portfolio Integration Server",
	}, nil)

	// Add get_portfolio tool
	readOnlyTrue := true
	openWorldFalse := false
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_portfolio",
		Title:       "Get Portfolio Balances",
		Description: "Retrieves current investment portfolio balances from KSEI (Indonesian Central Securities Depository) accounts. Returns detailed information about holdings including asset symbols, names, quantities, values, and currencies. Use this tool when you need to check current portfolio positions, asset allocations, or account balances. The data is fetched in real-time from KSEI AKSES and changes daily during settlement hours.",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Get Portfolio Balances",
			ReadOnlyHint:    true,
			IdempotentHint:  false, // Portfolio data changes over time (daily settlement updates)
			OpenWorldHint:   &openWorldFalse,
			DestructiveHint: &readOnlyTrue, // false means non-destructive (inverted from ReadOnly)
		},
	}, s.handleGetPortfolio)

	// Add list_account_names tool
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_account_names",
		Title:       "List Available Account Names",
		Description: "Lists all account names that are currently configured in the server. Use this tool to discover which accounts are available before calling get_portfolio with specific account names. Each account name represents a separate KSEI AKSES account connection. This is useful for understanding the scope of available data and for selecting specific accounts to query.",
		Annotations: &mcp.ToolAnnotations{
			Title:           "List Available Account Names",
			ReadOnlyHint:    true,
			IdempotentHint:  true,
			OpenWorldHint:   &openWorldFalse,
			DestructiveHint: &readOnlyTrue, // false means non-destructive
		},
	}, s.handleListAccountNames)

	s.mcpServer = mcpServer

	return s
}

// handleGetPortfolio handles the get_portfolio MCP tool
func (m *MCP) handleGetPortfolio(ctx context.Context, req *mcp.CallToolRequest, args GetPortfolioArgs) (*mcp.CallToolResult, GetPortfolioResult, error) {
	result := GetPortfolioResult{}

	clients := m.selectKseiClients(args.AccountNames)
	if len(clients) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Selected accounts not found, available accounts are " + strings.Join(m.getKseiClientNames(), ", "),
				},
			},
			IsError: true,
		}, result, nil
	}

	balances, err := getAllBalances(clients)
	if err != nil {
		return nil, result, err
	}

	if len(balances) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "No portfolio balances found for selected accounts",
				},
			},
			IsError: true,
		}, result, nil
	}

	result.Balances = balances

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: result.Description(),
			},
		},
	}, result, nil
}

func (m *MCP) handleListAccountNames(ctx context.Context, req *mcp.CallToolRequest, args ListAccountNamesArgs) (*mcp.CallToolResult, ListAccountNamesResult, error) {
	result := ListAccountNamesResult{
		AccountNames: m.getKseiClientNames(),
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: result.Description(),
			},
		},
	}, result, nil
}
