package game

import "testing"

func TestFaction(t *testing.T) {
	t.Run("Create Faction", func(t *testing.T) {
		f := Faction{
			UUID: "0c1094b4-31a5-4c0b-903a-0db5c874aeba",
			Name: "test",
		}
		if f.UUID != "0c1094b4-31a5-4c0b-903a-0db5c874aeba" {
			t.Errorf(
				"Expected UUID=0c1094b4-31a5-4c0b-903a-0db5c874aeba, got %s",
				f.UUID,
			)
		}

		if f.Name != "test" {
			t.Errorf("Expected Name=test, got %s", f.Name)
		}
	})
}
