package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/pomquest/internal/utils"
)

// DefaultDBPath returns the global default path for the pomquest SQLite database.
func DefaultDBPath() string {
	return utils.DBPathGlobal()
}

// AutoMigrateDB runs GORM AutoMigrate for all domain data models.
func AutoMigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
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
}

// EnsureInitialData ensures basic initial models exist in DB (first open / post-reset seeding).
func EnsureInitialData(db *gorm.DB) error {
	var partyCount int64
	if err := db.Model(&Party{}).Count(&partyCount).Error; err != nil {
		return err
	}
	if partyCount == 0 {
		initialParty := Party{
			Name:             "Adventuring Party",
			CurrentPartySize: 0,
			MaxPartySize:     1,
			Gold:             1,
			FoodLbs:          100,
			Diet:             "balanced",
			Pace:             "fast",
			SpareParts:       2,
			Oxen:             2,
		}
		if err := db.Create(&initialParty).Error; err != nil {
			return err
		}
	}

	var playerCount int64
	if err := db.Model(&Player{}).Count(&playerCount).Error; err != nil {
		return err
	}
	if playerCount == 0 {
		initialPlayer := Player{
			Name: "Dark Lord",
		}
		if err := db.Create(&initialPlayer).Error; err != nil {
			return err
		}
	}

	var hoardCount int64
	if err := db.Model(&Hoard{}).Count(&hoardCount).Error; err != nil {
		return err
	}
	if hoardCount == 0 {
		initialHoard := Hoard{
			GoldTotal: 1250,
		}
		if err := db.Create(&initialHoard).Error; err != nil {
			return err
		}
	}

	return nil
}

// LoadOrCreateActiveParty loads the active party from DB or initializes a default party if missing.
func LoadOrCreateActiveParty(db *gorm.DB) (Party, error) {
	if err := EnsureInitialData(db); err != nil {
		return Party{}, err
	}

	var party Party
	err := db.Preload("Heroes").First(&party).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			party = Party{
				Name:             "Adventuring Party",
				CurrentPartySize: 0,
				MaxPartySize:     1,
				Gold:             1,
			}
			if createErr := db.Create(&party).Error; createErr != nil {
				return party, createErr
			}
			return party, nil
		}
		return party, err
	}

	var activeHeroesCount int
	for _, h := range party.Heroes {
		if !h.IsDead {
			activeHeroesCount++
		}
	}
	party.CurrentPartySize = activeHeroesCount
	if party.MaxPartySize <= 0 {
		party.MaxPartySize = 1
	}
	return party, nil
}

// InitDB opens a SQLite database connection at dbPath, runs migrations, and seeds initial data.
func InitDB(dbPath string) (*gorm.DB, error) {
	if dbPath != ":memory:" && dbPath != "file::memory:?cache=shared" {
		if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.Exec("PRAGMA foreign_keys = ON;")
	db.Exec("PRAGMA journal_mode = WAL;")

	if err := AutoMigrateDB(db); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	if err := EnsureInitialData(db); err != nil {
		return nil, fmt.Errorf("failed to seed initial database records: %w", err)
	}

	return db, nil
}

// ResetAndMigrateDB resets the database by removing existing tables/file and auto-migrating all models.
func ResetAndMigrateDB(dbPath string) (*gorm.DB, error) {
	if dbPath != ":memory:" && dbPath != "file::memory:?cache=shared" {
		_ = os.Remove(dbPath)
		_ = os.Remove(dbPath + "-wal")
		_ = os.Remove(dbPath + "-shm")
	}

	db, err := InitDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to re-initialize database after reset: %w", err)
	}

	if dbPath == ":memory:" || dbPath == "file::memory:?cache=shared" {
		err = db.Migrator().DropTable(
			&Hero{}, &ClassDef{}, &Monster{}, &NPC{}, &Party{},
			&Hoard{}, &Quest{}, &Location{}, &Encounter{}, &GameState{},
			&Player{}, &GameLog{},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to drop tables for in-memory reset: %w", err)
		}
		if err := AutoMigrateDB(db); err != nil {
			return nil, fmt.Errorf("failed to auto-migrate after dropping tables: %w", err)
		}
		if err := EnsureInitialData(db); err != nil {
			return nil, fmt.Errorf("failed to re-seed initial database records: %w", err)
		}
	}

	return db, nil
}

// ReplaceOrAddHeroInParty saves a recruited hero into the active party in SQLite DB as passed.
func ReplaceOrAddHeroInParty(db *gorm.DB, index int, hero Hero) error {
	party, err := LoadOrCreateActiveParty(db)
	if err != nil {
		return err
	}

	hero.PartyID = &party.ID

	if index >= 0 && index < len(party.Heroes) {
		existingID := party.Heroes[index].ID
		hero.ID = existingID
		if err := db.Save(&hero).Error; err != nil {
			return err
		}
	} else {
		if err := db.Create(&hero).Error; err != nil {
			return err
		}
	}

	var activeCount int64
	db.Model(&Hero{}).Where("party_id = ? AND is_dead = ?", party.ID, false).Count(&activeCount)
	db.Model(&Party{}).Where("id = ?", party.ID).Update("current_party_size", activeCount)

	return nil
}

// LoadActivePartyFromPath opens/connects to the SQLite DB at dbPath and loads or creates the active party.
func LoadActivePartyFromPath(dbPath string) (Party, error) {
	if dbPath == "" {
		dbPath = DefaultDBPath()
	}
	db, err := InitDB(dbPath)
	if err != nil {
		return Party{
			Name:             "Adventuring Party",
			CurrentPartySize: 0,
			MaxPartySize:     1,
			Gold:             1,
		}, err
	}
	return LoadOrCreateActiveParty(db)
}

// SaveRecruitedHeroToPath opens/connects to the SQLite DB at dbPath and saves/replaces a hero in the active party.
func SaveRecruitedHeroToPath(dbPath string, targetIndex int, hero Hero) error {
	if dbPath == "" {
		dbPath = DefaultDBPath()
	}
	db, err := InitDB(dbPath)
	if err != nil {
		return err
	}
	return ReplaceOrAddHeroInParty(db, targetIndex, hero)
}

// SaveParty saves or updates party inventory supplies and coin totals in SQLite DB.
func SaveParty(db *gorm.DB, party Party) error {
	if party.ID == 0 {
		existing, err := LoadOrCreateActiveParty(db)
		if err == nil && existing.ID != 0 {
			party.ID = existing.ID
		}
	}
	return db.Save(&party).Error
}

// SavePartyToPath opens DB at dbPath and saves party inventory supplies and coins.
func SavePartyToPath(dbPath string, party Party) error {
	if dbPath == "" {
		dbPath = DefaultDBPath()
	}
	db, err := InitDB(dbPath)
	if err != nil {
		return err
	}
	return SaveParty(db, party)
}
