package game

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Position struct {
	X float64 `json:"X"`
	Y float64 `json:"Y"`
}

type Troop struct {
	UUID        string      `json:"UUID"`
	Position    *Position   `json:"Position"`
	Status      TroopStatus `json:"Status"`
	FactionUUID string      `json:"FactionUUID"`
	MatchUUID   string      `json:"MatchUUID"`
}

type TroopStatus int

const (
	Alive = iota
	Dead
	OffPlanet
)

func ParseTroopStatus(str string) (TroopStatus, error) {
	switch strings.ToLower(strings.TrimSpace(str)) {
	case "alive":
		return Alive, nil
	case "dead":
		return Dead, nil
	case "off planet", "offplanet":
		return OffPlanet, nil
	default:
		return 0, fmt.Errorf("invalid TroopStatus: %q", str)
	}
}

func (s TroopStatus) String() string {
	switch s {
	case Alive:
		return "Alive"
	case Dead:
		return "Dead"
	case OffPlanet:
		return "OffPlanet"
	default:
		return fmt.Sprintf("Unknown(%d)", s)
	}
}

func (s TroopStatus) MarshalJSON() ([]byte, error) {
	if !s.Valid() {
		return nil, fmt.Errorf("invalid TroopStatus: %d", s)
	}
	return []byte(`"` + s.String() + `"`), nil
}

func (s *TroopStatus) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		status, err := ParseTroopStatus(str)
		if err != nil {
			return err
		}
		*s = status
		return nil
	}

	// If string fails, try as number
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		if num >= Alive && num <= OffPlanet {
			*s = TroopStatus(num)
			return nil
		}
		return fmt.Errorf("invalid TroopStatus value: %d", num)
	}

	return fmt.Errorf("TroopStatus should be a string or number, got: %s", data)
}

func (s TroopStatus) Valid() bool {
	return s >= Alive && s <= OffPlanet
}
