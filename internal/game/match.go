package game

type UserInGame struct {
	UUID        string `json:"UUID"`
	MatchUUID   string `json:"MatchUUID"`
	FactionUUID string `json:"FactionUUID "`

	// cards
}

type Match struct {
	UUID         string            `json:"UUID"`
	GameStage    GameStage         `json:"GameStage"`
	FactionUsers map[string]string `json:"FactionUsers"`
}

type MatchStageType int

const (
	PickTraitor = iota
	StormMove
)

type GameStage struct {
	stage MatchStageType
	next  *GameStage
}

func (g GameStage) Next() GameStage {
	return *g.next
}
