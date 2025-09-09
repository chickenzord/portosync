package server

import (
	"fmt"
	"strings"
)

type GetPortfolioArgs struct {
	AccountNames []string `json:"account_names" jsonschema:"description:List of account names, empty means all"`
}

type Balance struct {
	SourceType    string  `json:"source_type"    jsonschema:"description:Type of the source"`
	SourceAccount string  `json:"source_account" jsonschema:"description:Name of the source account"`
	AssetSymbol   string  `json:"asset_symbol"   jsonschema:"description:Symbol of the asset"`
	AssetName     string  `json:"asset_name"     jsonschema:"description:Name of the asset"`
	AssetType     string  `json:"asset_type"     jsonschema:"description:Type of the asset"`
	UnitsAmount   float64 `json:"units_amount"   jsonschema:"description:Amount of the asset units"`
	UnitsValue    float64 `json:"units_value"    jsonschema:"description:Total value of the asset"`
	UnitsCurrency string  `json:"units_currency" jsonschema:"description:Currency of the asset units"`
}

func (b Balance) Description() string {
	return fmt.Sprintf("%s (%s): %f units, total value %s %f",
		b.AssetSymbol,
		b.SourceAccount,
		b.UnitsAmount,
		b.UnitsCurrency,
		b.UnitsValue,
	)
}

type GetPortfolioResult struct {
	Balances []Balance `json:"balances" jsonschema:"description:List of account balances"`
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
