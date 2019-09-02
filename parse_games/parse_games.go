package main

import (
	"github.com/lokken/family-football/gobbler"
	"github.com/lokken/family-football/types"
	"golang.org/x/net/html"
	"log"
	"strings"
	"sync"
)

func main() {
	rootNode := loadWeek()
	extractGames(rootNode)
}

func loadWeek() (rootNode *html.Node) {
	week := 1
	file := gobbler.LoadSchedule(week)

	rootNode, err := html.Parse(file)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func extractGames(rootNode *html.Node) {
	schedContainerNode := findSchedContainer(rootNode)
	// fmt.Printf("schedContainer: %s\n\n", schedContainerNode.Attr)

	var schedTables []*html.Node
	findSchedTables(&schedTables, schedContainerNode)
	// fmt.Printf("schedTables: %s\n\n", schedTables)

	var wg sync.WaitGroup
	gameNodeChan := make(chan *html.Node, 1)
	for _, tableNode := range schedTables {
		wg.Add(1)

		go func(tableNode *html.Node) {
			tbody := findTbody(tableNode)
			// fmt.Printf("tbody: %s\n\n", tbody)

			findSchedGames(tbody, gameNodeChan)
			wg.Done()
		}(tableNode)
	}

	go func() {
		wg.Wait()
		close(gameNodeChan)
	}()

	var wg2 sync.WaitGroup
	gameChan := make(chan types.Game, 1)
	for gameNode := range gameNodeChan {
		wg2.Add(1)

		go func(gameNode *html.Node) {
			parseGame(gameNode, gameChan)
			wg2.Done()
		}(gameNode)
	}

	go func() {
		wg2.Wait()
		close(gameChan)
	}()

	var games []types.Game
	for game := range gameChan {
		// fmt.Printf("game: %s\n\n", game)
		games = append(games, game)
	}
	gobbler.SaveGames(games)
}

func parseGame(n *html.Node, gameChan chan types.Game) {
	var gameTds []*html.Node
	parseGameHelper(&gameTds, n)

	game := types.Game{}

	for i, gameTd := range gameTds {
		switch i {
		case 0:
			game.AwayTeam = findTeam(gameTd)
		case 1:
			game.HomeTeam = findTeam(gameTd)
		case 2:
			for _, attr := range gameTd.Attr {
				if attr.Key == "data-date" {
					game.Time = strings.Replace(attr.Val, "Z", ":00Z", 1)
				}
			}
		case 5:
			if gameTd.FirstChild.Type == html.TextNode {
				location := gameTd.FirstChild.Data
				commaIndex := strings.Index(location, ",")
				if commaIndex > -1 {
					game.Stadium = strings.Trim(location[0:commaIndex], " ")
					game.Location = strings.Trim(location[commaIndex+1:], " ")
				}
			} else {
				for c := gameTd.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.ElementNode && c.Data == "a" {
						game.Stadium = strings.Trim(c.FirstChild.Data, " ")
					} else if c.Type == html.TextNode {
						game.Location = strings.Trim(c.Data, ", ")
					}
				}
			}
		}
	}

	gameChan <- game
}

func findTeam(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "abbr" {
		for _, attr := range n.Attr {
			if attr.Key == "title" {
				return attr.Val
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		found := findTeam(c)
		if len(found) > 0 {
			return found
		}
	}
	return ""
}

func parseGameHelper(found *[]*html.Node, n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "td" {
		*found = append(*found, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseGameHelper(found, c)
	}
}

func findSchedGames(n *html.Node, gameNodeChan chan *html.Node) {
	if n.Type == html.ElementNode && n.Data == "tr" {
		gameNodeChan <- n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findSchedGames(c, gameNodeChan)
	}
}

func findTbody(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "tbody" {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		found := findTbody(c)
		if found != nil {
			return found
		}
	}
	return nil
}

func findSchedTables(found *[]*html.Node, n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "table" {
		for _, attr := range n.Attr {
			if attr.Key == "class" {
				classes := strings.Split(attr.Val, " ")
				for _, class := range classes {
					if class == "schedule" {
						*found = append(*found, n)
					}
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findSchedTables(found, c)
	}
}

func findSchedContainer(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == "sched-container" {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		schedContainerNode := findSchedContainer(c)
		if schedContainerNode != nil {
			return schedContainerNode
		}
	}
	return nil
}
