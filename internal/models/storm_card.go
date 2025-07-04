package models

const MAX_STORM_STEPS = 6

type StormCard struct {
	UUID  string `json:"uuid"`
	steps uint
}

func NewStormCard(uuid string, steps uint) StormCard {
	return StormCard{
		UUID:  uuid,
		steps: steps,
	}
}

func (s *StormCard) Steps() uint {
	return s.steps
}
