package app

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAutoMigrateDataModels(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to memory database: %v", err)
	}

	err = db.AutoMigrate(
		&Hero{},
		&ClassDef{},
		&Monster{},
		&NPC{},
		&Party{},
		&Hoard{},
		&Quest{},
		&Location{},
		&Encounter{},
		&GameState{},
		&Player{},
		&GameLog{},
	)
	if err != nil {
		t.Fatalf("failed to auto migrate models: %v", err)
	}
}

func TestCreatureAndNounEmbedding(t *testing.T) {
	hero := Hero{
		Creature: Creature{
			Noun: Noun{
				Name: "Jimmothy",
				Icon: "⚔",
			},
			HP:     10,
			MaxHP:  10,
			Attack: 2,
			Armor:  2,
		},
		Class: "Fighter",
	}

	if hero.Name != "Jimmothy" {
		t.Errorf("expected hero.Name to be 'Jimmothy', got '%s'", hero.Name)
	}
	if hero.HP != 10 {
		t.Errorf("expected hero.HP to be 10, got %d", hero.HP)
	}
}
