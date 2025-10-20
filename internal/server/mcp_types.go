package server

import (
	"fmt"
	"strings"
)

type GetPortfolioArgs struct {
	AccountNames []string `json:"account_names" jsonschema:"description:List of specific account names to retrieve portfolio data from. Each name must match a configured account. If empty or omitted, returns portfolio data from all configured accounts. Use the list_account_names tool to discover available account names."`
}

type Balance struct {
	SourceType    string  `json:"source_type"    jsonschema:"description:Type of data source providing this balance (e.g., KSEI for Indonesian securities depository)"`
	SourceAccount string  `json:"source_account" jsonschema:"description:The account name from which this balance was retrieved, matching one of the configured account names"`
	AssetSymbol   string  `json:"asset_symbol"   jsonschema:"description:Trading symbol or ticker of the asset (e.g., BBCA for Bank Central Asia stock)"`
	AssetName     string  `json:"asset_name"     jsonschema:"description:Full descriptive name of the asset"`
	AssetType     string  `json:"asset_type"     jsonschema:"description:Primary classification of the asset (e.g., Stock, Bond, Mutual Fund)"`
	AssetSubType  string  `json:"asset_sub_type" jsonschema:"description:Additional classification or subtype of the asset, providing more granular categorization"`
	UnitsAmount   float64 `json:"units_amount"   jsonschema:"description:Quantity of asset units held in the account"`
	UnitsValue    float64 `json:"units_value"    jsonschema:"description:Total monetary value of the asset holdings in the specified currency"`
	UnitsCurrency string  `json:"units_currency" jsonschema:"description:Currency code for the asset value (e.g., IDR for Indonesian Rupiah, USD for US Dollar)"`
}

func (b Balance) AssetTypeFull() string {
	fragments := []string{}

	if b.AssetType != "" {
		fragments = append(fragments, b.AssetType)
	}

	if b.AssetSubType != "" {
		fragments = append(fragments, b.AssetSubType)
	}

	return strings.Join(fragments, "/")
}

func (b Balance) Description() string {
	return fmt.Sprintf("%s %s: %f units of %s, total value %s %f (%s)",
		b.AssetSymbol,
		b.AssetName,
		b.UnitsAmount,
		b.AssetTypeFull(),
		b.UnitsCurrency,
		b.UnitsValue,
		b.SourceAccount,
	)
}

type GetPortfolioResult struct {
	Balances []Balance `json:"balances" jsonschema:"description:Array of portfolio balances across all requested accounts. Each balance represents a single asset holding with quantity and value information."`
}

// Description returns a description of the GetPortfolioResult as MCP response text
func (r GetPortfolioResult) Description() string {
	if len(r.Balances) == 0 {
		return "Portfolio is empty"
	}

	var descriptions []string
	for _, balance := range r.Balances {
		descriptions = append(descriptions, fmt.Sprintf("- %s", balance.Description()))
	}

	return "Portfolio:\n" + strings.Join(descriptions, "\n")
}

type ListAccountNamesArgs struct {
	// This tool requires no parameters - it returns all configured account names
}

type ListAccountNamesResult struct {
	AccountNames []string `json:"account_names" jsonschema:"description:Array of configured account names. These names can be used as parameters when calling the get_portfolio tool to filter results by specific accounts."`
}

// Description returns a description of the ListAccountNamesResult as MCP response text
func (r ListAccountNamesResult) Description() string {
	if len(r.AccountNames) == 0 {
		return "No accounts configured"
	}

	return fmt.Sprintf("Available accounts: %s", strings.Join(r.AccountNames, ", "))
}
