package main

import (
	//"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	//"os"
)

type Site struct {
	Title string
}

type ChooseGamesPage struct {
	Title string
	CancelButtonText string
	NextButtonText string
	Cards []GameCard
}

type ChooseBonusPage struct {
	Title string
	CancelButtonText string
	NextButtonText string
	Cards []BonusCard
}

type GameCard struct {
	AwayTeam string
	HomeTeam string
	EventTime time.Time
}

type BonusCard struct {
	Qualifier string
	Quantifier int32
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("out.html")
	if err != nil {
		panic(err)
	}

	t, err := time.Parse(time.RFC3339, "2019-08-24T23:00:00Z")
	if err != nil {
		panic(err)
	}

	page := struct {
		Site Site
		Page ChooseGamesPage
	}{
		Site: Site{Title: "Family Football 2019-20"},
		Page: ChooseGamesPage{
			Title: "Choose Games",
			CancelButtonText: "Cancel",
			NextButtonText: "Choose Bonus >>",
			Cards: []GameCard{
				{
					AwayTeam: "Texas Tech Red	Raiders",
					HomeTeam: "Texas Longhorns",
					EventTime: t,
				},
				{
					AwayTeam: "Oklahoma Sooners",
					HomeTeam: "West Virginia Mountaineers",
					EventTime: t,
				},
			},
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
