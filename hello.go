package main

import (
	"fmt"
	"github.com/lokken/family-football/gobbler"
	"github.com/lokken/family-football/types"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Site struct {
	Title string
}

type ChooseGamesPage struct {
	Title            string
	CancelButtonText string
	NextButtonText   string
	Cards            []GameCard
}

type ChooseBonusPage struct {
	Title            string
	CancelButtonText string
	NextButtonText   string
	Cards            []BonusCard
}

type GameCard struct {
	AwayTeam  string
	HomeTeam  string
	Stadium   string
	Location  string
	EventTime time.Time
}

type BonusCard struct {
	Qualifier  string
	Quantifier int32
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("out.html")
	if err != nil {
		panic(err)
	}

	var gameCards []GameCard
	var games []types.Game
	gobbler.LoadGames(&games)
	for _, game := range games {
		fmt.Printf("time: %s\n", game.Time)

		t, err := time.Parse(time.RFC3339, "1979-09-05T00:00:00Z")
		if err != nil {
			panic(err)
		}
		if len(game.Time) > 0 {
			t, err = time.Parse(time.RFC3339, game.Time)
			if err != nil {
				panic(err)
			}
		}
		fmt.Printf("t: %s\n", t)
		gameCard := GameCard{
			AwayTeam:  game.AwayTeam,
			HomeTeam:  game.HomeTeam,
			Stadium: game.Stadium,
			Location: game.Location,
			EventTime: t,
		}
		gameCards = append(gameCards, gameCard)
	}
	fmt.Printf("games: %s\n", games)

	page := struct {
		Site Site
		Page ChooseGamesPage
	}{
		Site: Site{Title: "Family Football 2019-20"},
		Page: ChooseGamesPage{
			Title:            "Choose Games",
			CancelButtonText: "Cancel",
			NextButtonText:   "Choose Bonus >>",
			Cards:            gameCards,
			// Cards: []GameCard{
			// 	{
			// 		AwayTeam: "Texas Tech Red	Raiders",
			// 		HomeTeam: "Texas Longhorns",
			// 		EventTime: t,
			// 	},
			// 	{
			// 		AwayTeam: "Oklahoma Sooners",
			// 		HomeTeam: "West Virginia Mountaineers",
			// 		EventTime: t,
			// 	},
			// },
		},
	}

	err = tmpl.Execute(w, page)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
