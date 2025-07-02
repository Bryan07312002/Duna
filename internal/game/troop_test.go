package game

import (
	"encoding/json"
	"strings"
	"testing"
)

// Removes Tabs, Spaces and line breaks
func compactJSON(s string) string {
	noTab := strings.ReplaceAll(s, "\t", "")
	noSpaces := strings.ReplaceAll(noTab, " ", "")
	return strings.ReplaceAll(noSpaces, "\n", "")
}

func TestPosition(t *testing.T) {
	t.Run("Create Position", func(t *testing.T) {
		p := Position{X: 10.5, Y: 20.5}
		if p.X != 10.5 {
			t.Errorf("Expected X=10.5, got %f", p.X)
		}

		if p.Y != 20.5 {
			t.Errorf("Expected Y=20.5, got %f", p.Y)
		}
	})
}

func TestTroop(t *testing.T) {
	t.Run("Create Troop", func(t *testing.T) {
		pos := &Position{X: 1.0}
		troop := Troop{
			UUID:     "123",
			Position: pos,
			Status:   Alive,
		}

		if troop.UUID != "123" {
			t.Errorf("Expected UUID=123, got %s", troop.UUID)
		}
		if troop.Position.X != 1.0 {
			t.Errorf("Expected Position.X=1.0, got %f", troop.Position.X)
		}
		if troop.Status != Alive {
			t.Errorf("Expected Status=Alive, got %v", troop.Status)
		}
	})
}

func TestTroopStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected TroopStatus
		wantErr  bool
	}{
		{"Alive lowercase", "alive", Alive, false},
		{"Alive mixed case", "AlIvE", Alive, false},
		{"Dead", "dead", Dead, false},
		{"OffPlanet with space", "off planet", OffPlanet, false},
		{"Invalid status", "unknown", 0, true},
		{"Empty string", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTroopStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTroopStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ParseTroopStatus() = %v, want %v", got, tt.expected)
			}
		})
	}

	t.Run("String() method", func(t *testing.T) {
		tests := []struct {
			status   TroopStatus
			expected string
			wantErr  bool
		}{
			{Alive, "Alive", false},
			{Dead, "Dead", false},
			{OffPlanet, "OffPlanet", false},
			{TroopStatus(99), "Unknown(99)", true},
		}

		for _, tt := range tests {
			got := tt.status.String()

			if got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		}
	})

	t.Run("Valid() method", func(t *testing.T) {
		tests := []struct {
			status   TroopStatus
			expected bool
		}{
			{Alive, true},
			{Dead, true},
			{OffPlanet, true},
			{-1, false},
			{3, false},
		}

		for _, tt := range tests {
			if got := tt.status.Valid(); got != tt.expected {
				t.Errorf("Valid() = %v, want %v for status %v", got, tt.expected, tt.status)
			}
		}
	})
}

func TestTroopStatusJSON(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		tests := []struct {
			status   TroopStatus
			expected string
			wantErr  bool
		}{
			{Alive, `"Alive"`, false},
			{Dead, `"Dead"`, false},
			{OffPlanet, `"OffPlanet"`, false},
		}

		for _, tt := range tests {
			got, err := json.Marshal(tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				continue
			}

			if string(got) != tt.expected {
				t.Errorf("MarshalJSON() = %v, want %v", string(got), tt.expected)
			}
		}
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		tests := []struct {
			input    string
			expected TroopStatus
			wantErr  bool
		}{
			{`"Alive"`, Alive, false},
			{`"alive"`, Alive, false},
			{`"Dead"`, Dead, false},
			{`"off planet"`, OffPlanet, false},
			{`"invalid"`, 0, true},
			{`123`, 0, true},
		}

		for _, tt := range tests {
			var s TroopStatus
			err := json.Unmarshal([]byte(tt.input), &s)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				continue
			}
			if s != tt.expected {
				t.Errorf("UnmarshalJSON() = %v, want %v", s, tt.expected)
			}
		}
	})
}

func TestTroopJSON(t *testing.T) {
	t.Run("Full Troop Marshaling", func(t *testing.T) {
		troop := Troop{
			UUID:        "7d8e9f53-b50e-4718-9a65-762350e09466",
			Position:    &Position{X: 1.5},
			Status:      OffPlanet,
			FactionUUID: "d6617fe9-4d3f-49d8-8edb-115f8ee44bbd",
			MatchUUID:   "760cd173-9c0c-4493-b50a-0c11ade901fe",
		}

		data, err := json.Marshal(troop)
		if err != nil {
			t.Fatalf("Failed to marshal troop: %v", err)
		}

		expected := `{
			"UUID": "7d8e9f53-b50e-4718-9a65-762350e09466",
			"Position": {
				"X": 1.5,
				"Y": 0
			},
			"Status": "OffPlanet",
			"FactionUUID": "d6617fe9-4d3f-49d8-8edb-115f8ee44bbd",
			"MatchUUID": "760cd173-9c0c-4493-b50a-0c11ade901fe"
		}`
		if string(data) != compactJSON(expected) {
			t.Errorf("Expected %s, got %s", compactJSON(expected), string(data))
		}

		var decoded Troop
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal troop: %v", err)
		}

		if decoded.UUID != troop.UUID || decoded.Status != troop.Status {
			t.Errorf("Decoded troop doesn't match original")
		}
	})
}
