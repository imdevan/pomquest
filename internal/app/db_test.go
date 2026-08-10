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
