package server

import (
	"sync"

	"github.com/chickenzord/goksei"
	"golang.org/x/sync/errgroup"
)

var (
	allPortfolioTypes = []goksei.PortfolioType{
		goksei.EquityType,
		goksei.BondType,
		goksei.MutualFundType,
	}
)

// getAllBalances retrieves all share balances from the client
// in parallel using map-reduce pattern
func getAllBalances(clients map[string]*goksei.Client) ([]Balance, error) {
	var mu sync.Mutex

	var errs errgroup.Group

	var balances []Balance

	for accountName, client := range clients {
		for _, portfolioType := range allPortfolioTypes {
			errs.Go(func() error {
				res, err := client.GetShareBalances(portfolioType)
				if err != nil {
					return err
				}

				mu.Lock()

				for _, b := range res.Data {
					balances = append(balances, Balance{
						SourceType:    "ksei",
						SourceAccount: accountName,
						AssetSymbol:   b.Symbol(),
						AssetName:     b.Name(),
						AssetType:     portfolioType.Name(),
						UnitsCurrency: b.Currency,
						UnitsAmount:   b.Amount,
						UnitsValue:    b.CurrentValue(),
					})
				}

				mu.Unlock()

				return nil
			})
		}
	}

	if err := errs.Wait(); err != nil {
		return nil, err
	}

	return balances, nil
}
