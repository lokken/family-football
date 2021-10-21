package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/lokken/family-football/gobbler"
	"github.com/lokken/family-football/types"
)

func main() {
	// http.HandleFunc("/week/", weekHandler)
	// http.HandleFunc("/leaderboard/", leaderboardHandler)
	// http.HandleFunc("/player/", playerHandler)

	// http.HandleFunc("/choose-games/", chooseGamesHandler)
	// http.HandleFunc("/choose-bonus/", chooseBonusHandler)
	http.HandleFunc("/configure-bonus/", configureBonusHandler)
	http.HandleFunc("/configure-bonus", configureBonusHandler)

	http.HandleFunc("/delete-bonus/", deleteBonusHandler)
	http.HandleFunc("/delete-bonus", deleteBonusHandler)

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
	Cards            map[string][]types.Bonus
}

type configureBonusPage struct {
	Title string
	Cards map[string][]types.Bonus
}

func tmpl(which string) string {
	_, filePath, _, _ := runtime.Caller(0)
	filename := fmt.Sprintf("%s.html", which)
	return path.Join(filepath.Dir(filePath), "tmpl", filename)
}

func weekHandler(w http.ResponseWriter, r *http.Request) {
	// path := r.URL.Path[len("/week/"):]
}

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
}

func playerHandler(w http.ResponseWriter, r *http.Request) {
}

func deleteBonusHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("%s\n", r.Form)
}

func configureBonusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		command := strings.ToLower(r.Header.Get("X-Command"))
		if command == "" {
			http.Error(w, "Bad X-Command", http.StatusMultipleChoices)
			return
		}

		var data []types.Bonus
		body, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(body, &data)

		var bonuses map[string]*types.Bonus
		gobbler.LoadBonuses(&bonuses)

		fmt.Printf("%s\n", body)
		fmt.Printf("%s\n", data)
		switch command {
		case "save":
			for _, b := range data {
				bonuses[b.ID].Qualifier = b.Qualifier
				bonuses[b.ID].Quantifier = b.Quantifier
			}
			gobbler.PutBonuses(bonuses)
		case "delete":
			for _, b := range data {
				delete(bonuses, b.ID)
			}
			gobbler.SaveBonuses(bonuses)
		}

		w.WriteHeader(http.StatusNoContent)
		w.Header().Set("Content-Type", "application/json")

	} else if r.Method == "GET" {
		tmpl, _ := template.ParseFiles(tmpl("configure-bonus-formpost"))

		var bonuses map[string]*types.Bonus
		gobbler.LoadBonuses(&bonuses)

		bonusMap := make(map[string][]types.Bonus)
		for id, b := range bonuses {
			b.ID = id
			bonusMap[b.Type] = append(bonusMap[b.Type], *b)
		}

		page := struct {
			Site site
			Page configureBonusPage
		}{
			Site: site{Title: "Family Football 2019-20", Copywrite: "© Keith Lokken"},
			Page: configureBonusPage{
				Title: "Configure Bonuses",
				Cards: bonusMap,
			},
		}

		_ = tmpl.Execute(w, page)
	} else {
		http.Error(w, "Invalid", http.StatusMethodNotAllowed)
		return
	}
}

func chooseBonusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid", http.StatusMethodNotAllowed)
		return
	}

	tmpl, _ := template.ParseFiles(tmpl("choose-bonus"))

	var bonuses map[string]*types.Bonus
	gobbler.LoadBonuses(&bonuses)

	bonusMap := make(map[string][]types.Bonus)
	for id, b := range bonuses {
		b.ID = id
		bonusMap[b.Type] = append(bonusMap[b.Type], *b)
	}

	page := struct {
		Site site
		Page chooseBonusPage
	}{
		Site: site{Title: "Family Football 2019-20", Copywrite: "© Keith Lokken"},
		Page: chooseBonusPage{
			Title:            "Choose Bonus",
			CancelButtonText: "Cancel",
			NextButtonText:   "Complete",
			Cards:            bonusMap,
		},
	}

	_ = tmpl.Execute(w, page)
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
