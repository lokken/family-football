package types

type Game struct {
	ID       string
	AwayTeam string
	HomeTeam string
	Time     string
	Stadium  string
	Location string
}

type Bonus struct {
	ID         string
	Type       string
	Qualifier  string
	Quantifier string
}
