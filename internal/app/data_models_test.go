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

	party := Party{
		Name:      "Fellowship",
		Loot:      5,
		Rope:      2,
		Bombs:     3,
		Potions:   4,
		Shovels:   1,
		ItemUsage: "sparse",
		Greed:     "low",
		Pace:      "fast",
		Diet:      "balanced",
	}

	if err := db.Create(&party).Error; err != nil {
		t.Fatalf("failed to create party record: %v", err)
	}

	var fetched Party
	if err := db.First(&fetched, party.ID).Error; err != nil {
		t.Fatalf("failed to fetch party record: %v", err)
	}

	if fetched.Loot != 5 || fetched.Rope != 2 || fetched.Bombs != 3 || fetched.Shovels != 1 {
		t.Errorf("unexpected supply quantities: %+v", fetched)
	}
	if fetched.ItemUsage != "sparse" || fetched.Greed != "low" {
		t.Errorf("unexpected item usage or greed: %+v", fetched)
	}
	if fetched.CurrentPartySize != 0 {
		t.Errorf("expected default CurrentPartySize to be 0, got %d", fetched.CurrentPartySize)
	}
	if fetched.MaxPartySize != 1 {
		t.Errorf("expected default MaxPartySize to be 1, got %d", fetched.MaxPartySize)
	}

	customParty := Party{
		Name:             "Custom Party",
		CurrentPartySize: 3,
		MaxPartySize:     5,
	}
	if err := db.Create(&customParty).Error; err != nil {
		t.Fatalf("failed to create custom party record: %v", err)
	}

	var fetchedCustom Party
	if err := db.First(&fetchedCustom, customParty.ID).Error; err != nil {
		t.Fatalf("failed to fetch custom party record: %v", err)
	}
	if fetchedCustom.CurrentPartySize != 3 || fetchedCustom.MaxPartySize != 5 {
		t.Errorf("expected party size 3/5, got %d/%d", fetchedCustom.CurrentPartySize, fetchedCustom.MaxPartySize)
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
