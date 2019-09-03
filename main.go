package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/lokken/family-football/gobbler"
	"github.com/lokken/family-football/types"
)

func main() {
	http.HandleFunc("/choose-games/", chooseGamesHandler)
	http.HandleFunc("/choose-bonus/", chooseBonusHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type site struct {
	Title     string
	Copywrite string
}

type chooseGamesPage struct {
	Title            string
	CancelButtonText string
	NextButtonText   string
	Cards            []types.Game
}

type chooseBonusPage struct {
	Title            string
	CancelButtonText string
	NextButtonText   string
	Cards            []types.Bonus
}

func tmpl(which string) string {
	_, filePath, _, _ := runtime.Caller(0)
	filename := fmt.Sprintf("%s.html", which)
	return path.Join(filepath.Dir(filePath), "tmpl", filename)
}

func chooseBonusHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(tmpl("choose-bonus"))
	if err != nil {
		panic(err)
	}

	var bonuses []types.Bonus
	gobbler.LoadBonuses(&bonuses)

	page := struct {
		Site site
		Page chooseBonusPage
	}{
		Site: site{Title: "Family Football 2019-20", Copywrite: "© Keith Lokken"},
		Page: chooseBonusPage{
			Title:            "Choose Bonus",
			CancelButtonText: "Cancel",
			NextButtonText:   "Complete",
			Cards:            bonuses,
		},
	}

	err = tmpl.Execute(w, page)
	if err != nil {
		panic(err)
	}
}

func chooseGamesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(tmpl("choose-games"))
	if err != nil {
		panic(err)
	}

	var games []types.Game
	gobbler.LoadGames(&games)
	loc, _ := time.LoadLocation("Local")
	for i := range games {
		game := &games[i]
		// fmt.Printf("time: %s\n", game.Time)

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
		game.Time = t.In(loc).Format("Monday, Jan _2 2006 3:04 PM MST")
	}
	// fmt.Printf("games: %s\n", games)

	page := struct {
		Site site
		Page chooseGamesPage
	}{
		Site: site{Title: "Family Football 2019-20", Copywrite: "© Keith Lokken"},
		Page: chooseGamesPage{
			Title:            "Choose Games",
			CancelButtonText: "Cancel",
			NextButtonText:   "Choose Bonus >>",
			Cards:            games,
		},
	}

	err = tmpl.Execute(w, page)
	if err != nil {
		panic(err)
	}
}
