package main

import (
	"github.com/lokken/family-football/gobbler"
	"github.com/lokken/family-football/types"
)

func main() {
	bonuses := []types.Bonus{
		{
			Type:       "Passing",
			Qualifier:  "What *team* earned the most passing yards?",
			Quantifier: "How many passing yards did they earn?",
		},
		{
			Type:       "Passing",
			Qualifier:  "What *quarterback* earned the most passing yards?",
			Quantifier: "How many passing yards did they earn?",
		},
		{
			Type:       "Rushing",
			Qualifier:  "What *team* earned the most rushing yards?",
			Quantifier: "How many rushing yards did they earn?",
		},
		{
			Type:       "Defense",
			Qualifier:  "What *player* earned the most tackles?",
			Quantifier: "How many tackles did they earn?",
		},
	}

	gobbler.SaveBonuses(bonuses)
}
