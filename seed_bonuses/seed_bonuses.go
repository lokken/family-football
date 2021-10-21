package main

import (
	"github.com/lokken/family-football/gobbler"
	"github.com/lokken/family-football/types"
	"github.com/rs/xid"
)

func main() {
	firstID := xid.New().String()
	bonuses := map[string]*types.Bonus{
		firstID: {
			Type:       "Passing",
			Qualifier:  "What *team* earned the most passing yards?",
			Quantifier: "How many passing yards did they earn?",
		},
		xid.New().String(): {
			Type:       "Passing",
			Qualifier:  "What *quarterback* earned the most passing yards?",
			Quantifier: "How many passing yards did they earn?",
		},
		xid.New().String(): {
			Type:       "Rushing",
			Qualifier:  "What *team* earned the most rushing yards?",
			Quantifier: "How many rushing yards did they earn?",
		},
		xid.New().String(): {
			Type:       "Defense",
			Qualifier:  "What *player* earned the most tackles?",
			Quantifier: "How many tackles did they earn?",
		},
	}

	gobbler.SaveBonuses(bonuses)

	bonusi := map[string]*types.Bonus{
		firstID: {
			Type:       "Passing",
			Qualifier:  "What team earned the most passing yards?",
			Quantifier: "How many passing yards did they earn?",
		},
		xid.New().String(): {
			Type:       "Defense",
			Qualifier:  "What *team* earned the most sacks?",
			Quantifier: "How many sacks did they earn?",
		},
	}
	gobbler.PutBonuses(bonusi)
}
