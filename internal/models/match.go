package models

type MatchState int

const (
	WaitingPlayers = iota
	InGame
	Finish
)

type Match struct {
	UUID          string
	MatchState    MatchState
	CreatedByUser string
}

func NewMatch(uuid string, matchState MatchState) Match {
	return Match{
		UUID:       uuid,
		MatchState: matchState,
	}
}
