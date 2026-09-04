package app

import (
	"path/filepath"
	"testing"
)

func TestInitAndResetDB(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_pomquest.db")

	db, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	// InitDB seeds 1 default party
	var count int64
	db.Model(&Party{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 initial seeded party, got %d", count)
	}

	// Insert custom record
	p := Party{Name: "Reset Test Party"}
	if err := db.Create(&p).Error; err != nil {
		t.Fatalf("failed to insert party: %v", err)
	}

	db.Model(&Party{}).Count(&count)
	if count != 2 {
		t.Fatalf("expected 2 parties, got %d", count)
	}

	// Reset and migrate (re-seeds 1 initial default party)
	resetDB, err := ResetAndMigrateDB(dbPath)
	if err != nil {
		t.Fatalf("ResetAndMigrateDB failed: %v", err)
	}

	resetDB.Model(&Party{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 initial party after reset, got %d", count)
	}
}

func TestResetAndMigrateDB_InMemory(t *testing.T) {
	db, err := InitDB("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("InitDB in-memory failed: %v", err)
	}

	p := Party{Name: "Memory Party"}
	if err := db.Create(&p).Error; err != nil {
		t.Fatalf("failed to insert party: %v", err)
	}

	resetDB, err := ResetAndMigrateDB("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("ResetAndMigrateDB in-memory failed: %v", err)
	}

	var count int64
	resetDB.Model(&Party{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 initial party after in-memory reset, got %d", count)
	}
}

func TestReplaceOrAddHeroInParty(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "recruitment_test.db")

	db, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	party, err := LoadOrCreateActiveParty(db)
	if err != nil {
		t.Fatalf("LoadOrCreateActiveParty failed: %v", err)
	}
	if party.CurrentPartySize != 0 {
		t.Fatalf("expected initial party size 0, got %d", party.CurrentPartySize)
	}

	hero := Hero{
		Creature: Creature{
			Noun:  Noun{Name: "Conan"},
			HP:    12,
			MaxHP: 12,
		},
		Class: "barbarian",
	}

	if err := ReplaceOrAddHeroInParty(db, 0, hero); err != nil {
		t.Fatalf("ReplaceOrAddHeroInParty failed: %v", err)
	}

	updatedParty, err := LoadOrCreateActiveParty(db)
	if err != nil {
		t.Fatalf("failed to reload party: %v", err)
	}
	if updatedParty.CurrentPartySize != 1 {
		t.Fatalf("expected updated party size 1, got %d", updatedParty.CurrentPartySize)
	}
	if len(updatedParty.Heroes) != 1 || updatedParty.Heroes[0].Name != "Conan" {
		t.Fatalf("expected hero Conan in party, got %+v", updatedParty.Heroes)
	}
}

func TestSaveSelectedLocationAndLoadGameState(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "location_test.db")

	loc := Location{
		Name:               "The Sunken Tomb of Eldoria",
		Biome:              "Ancient Catacombs",
		Description:        "A dark tomb sunken beneath mossy swamps.",
		DifficultyModifier: 1.2,
	}

	savedLoc, err := SaveSelectedLocationToPath(dbPath, loc)
	if err != nil {
		t.Fatalf("SaveSelectedLocationToPath failed: %v", err)
	}
	if savedLoc.ID == 0 {
		t.Fatalf("expected non-zero ID for saved location")
	}

	gameState, err := LoadActiveGameStateFromPath(dbPath)
	if err != nil {
		t.Fatalf("LoadActiveGameStateFromPath failed: %v", err)
	}

	if gameState.CurrentLocationID == nil || *gameState.CurrentLocationID != savedLoc.ID {
		t.Fatalf("expected CurrentLocationID %d, got %v", savedLoc.ID, gameState.CurrentLocationID)
	}
	if gameState.CurrentLocation == nil || gameState.CurrentLocation.Name != "The Sunken Tomb of Eldoria" {
		t.Fatalf("expected location name 'The Sunken Tomb of Eldoria', got %+v", gameState.CurrentLocation)
	}
}

func TestIsPartyQuesting(t *testing.T) {
	partyWithLiving := Party{
		Heroes: []Hero{
			{Creature: Creature{Noun: Noun{Name: "Hero1"}, HP: 10, MaxHP: 10}},
		},
	}
	partyDead := Party{
		Heroes: []Hero{
			{Creature: Creature{Noun: Noun{Name: "Hero1"}, HP: 0, MaxHP: 10}, IsDead: true},
		},
	}
	emptyParty := Party{Heroes: []Hero{}}

	questActive := &Quest{
		Status:       "Active",
		DungeonState: "InRoom",
	}
	questCompleted := &Quest{
		Status:       "Completed",
		DungeonState: "Completed",
	}
	questFailed := &Quest{
		Status:       "Failed",
		DungeonState: "Defeat",
	}

	if !IsPartyQuesting(partyWithLiving, GameState{CurrentQuest: questActive}) {
		t.Errorf("expected IsPartyQuesting to be true for living party and active quest")
	}
	if IsPartyQuesting(partyDead, GameState{CurrentQuest: questActive}) {
		t.Errorf("expected IsPartyQuesting to be false for dead party")
	}
	if IsPartyQuesting(emptyParty, GameState{CurrentQuest: questActive}) {
		t.Errorf("expected IsPartyQuesting to be false for empty party")
	}
	if IsPartyQuesting(partyWithLiving, GameState{CurrentQuest: questCompleted}) {
		t.Errorf("expected IsPartyQuesting to be false for completed quest")
	}
	if IsPartyQuesting(partyWithLiving, GameState{CurrentQuest: questFailed}) {
		t.Errorf("expected IsPartyQuesting to be false for failed quest")
	}
	if IsPartyQuesting(partyWithLiving, GameState{CurrentQuest: nil}) {
		t.Errorf("expected IsPartyQuesting to be false for nil quest")
	}
}


