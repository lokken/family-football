package types

type Game struct {
	AwayTeam string
	HomeTeam string
	Time     string
	Stadium  string
	Location string
}

type Bonus struct {
	Qualifier  string
	Quantifier string
	Type       string
}
