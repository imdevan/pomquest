package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

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
			Coins:            1000,
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
				Coins:            1000,
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

var (
	dbCacheMu sync.Mutex
	dbCache   = make(map[string]*gorm.DB)
)

// CloseDB closes and removes a cached DB connection.
func CloseDB(dbPath string) error {
	dbCacheMu.Lock()
	defer dbCacheMu.Unlock()
	if db, exists := dbCache[dbPath]; exists {
		delete(dbCache, dbPath)
		sqlDB, err := db.DB()
		if err == nil && sqlDB != nil {
			return sqlDB.Close()
		}
	}
	return nil
}

// CloseAllDBs closes all cached database connections.
func CloseAllDBs() {
	dbCacheMu.Lock()
	defer dbCacheMu.Unlock()
	for p, db := range dbCache {
		delete(dbCache, p)
		sqlDB, err := db.DB()
		if err == nil && sqlDB != nil {
			_ = sqlDB.Close()
		}
	}
}

func initDBUncached(dbPath string) (*gorm.DB, error) {
	if dbPath != ":memory:" && dbPath != "file::memory:?cache=shared" {
		if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Discard,
	})
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

// InitDB returns a cached SQLite database connection at dbPath or initializes and caches one.
func InitDB(dbPath string) (*gorm.DB, error) {
	dbCacheMu.Lock()
	defer dbCacheMu.Unlock()

	if db, exists := dbCache[dbPath]; exists {
		sqlDB, err := db.DB()
		if err == nil && sqlDB.Ping() == nil {
			return db, nil
		}
		delete(dbCache, dbPath)
	}

	db, err := initDBUncached(dbPath)
	if err != nil {
		return nil, err
	}
	dbCache[dbPath] = db
	return db, nil
}

// ResetAndMigrateDB resets the database by removing existing tables/file and auto-migrating all models.
func ResetAndMigrateDB(dbPath string) (*gorm.DB, error) {
	_ = CloseDB(dbPath)

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
			Coins:            1000,
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

// SaveSelectedLocation creates the location in DB and updates active GameState.
func SaveSelectedLocation(db *gorm.DB, loc Location) (Location, error) {
	if err := db.Create(&loc).Error; err != nil {
		return loc, err
	}

	var state GameState
	err := db.Preload("CurrentLocation").Where("is_active = ?", true).First(&state).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			state = GameState{
				IsActive: true,
			}
			if createErr := db.Create(&state).Error; createErr != nil {
				return loc, createErr
			}
		} else {
			return loc, err
		}
	}

	state.CurrentLocationID = &loc.ID
	state.CurrentLocation = &loc
	state.CurrentQuestID = nil
	state.CurrentQuest = nil
	if err := db.Save(&state).Error; err != nil {
		return loc, err
	}

	_ = ClearGameLogs(db, state.ID)

	return loc, nil
}

// SaveSelectedLocationToPath opens DB at dbPath and saves location to active GameState.
func SaveSelectedLocationToPath(dbPath string, loc Location) (Location, error) {
	if dbPath == "" {
		dbPath = DefaultDBPath()
	}
	db, err := InitDB(dbPath)
	if err != nil {
		return loc, err
	}
	return SaveSelectedLocation(db, loc)
}

// LoadActiveGameState loads active session GameState from DB.
func LoadActiveGameState(db *gorm.DB) (GameState, error) {
	var state GameState
	err := db.Preload("CurrentLocation").Preload("Party").Preload("Hoard").Preload("CurrentQuest").Where("is_active = ?", true).First(&state).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			state = GameState{
				IsActive: true,
			}
			if createErr := db.Create(&state).Error; createErr != nil {
				return state, createErr
			}
			return state, nil
		}
		return state, err
	}
	return state, nil
}

// LoadActiveGameStateFromPath opens DB at dbPath and loads active GameState.
func LoadActiveGameStateFromPath(dbPath string) (GameState, error) {
	if dbPath == "" {
		dbPath = DefaultDBPath()
	}
	db, err := InitDB(dbPath)
	if err != nil {
		return GameState{}, err
	}
	return LoadActiveGameState(db)
}

// SaveActiveQuest creates or updates the active quest in DB and sets it on active GameState.
func SaveActiveQuest(db *gorm.DB, quest Quest) (Quest, error) {
	if quest.ID == 0 {
		if err := db.Create(&quest).Error; err != nil {
			return quest, err
		}
	} else {
		if err := db.Save(&quest).Error; err != nil {
			return quest, err
		}
	}

	var state GameState
	err := db.Where("is_active = ?", true).First(&state).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			state = GameState{
				IsActive: true,
			}
			if createErr := db.Create(&state).Error; createErr != nil {
				return quest, createErr
			}
		} else {
			return quest, err
		}
	}

	state.CurrentQuestID = &quest.ID
	state.CurrentQuest = &quest
	if err := db.Save(&state).Error; err != nil {
		return quest, err
	}

	return quest, nil
}

// SaveActiveQuestToPath opens DB at dbPath and saves active quest state.
func SaveActiveQuestToPath(dbPath string, quest Quest) (Quest, error) {
	if dbPath == "" {
		dbPath = DefaultDBPath()
	}
	db, err := InitDB(dbPath)
	if err != nil {
		return quest, err
	}
	return SaveActiveQuest(db, quest)
}

// SaveGameLogs replaces/persists log messages for the given gameStateID.
func SaveGameLogs(db *gorm.DB, gameStateID uint, messages []string) error {
	if gameStateID == 0 {
		return nil
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("game_state_id = ?", gameStateID).Delete(&GameLog{}).Error; err != nil {
			return err
		}
		if len(messages) == 0 {
			return nil
		}
		gLogs := make([]GameLog, len(messages))
		for i, msg := range messages {
			gLogs[i] = GameLog{
				GameStateID: gameStateID,
				Message:     msg,
				Category:    "dungeon",
			}
		}
		return tx.Create(&gLogs).Error
	})
}

// SaveGameLogsToPath opens DB at dbPath and saves log messages.
func SaveGameLogsToPath(dbPath string, gameStateID uint, messages []string) error {
	if dbPath == "" {
		dbPath = DefaultDBPath()
	}
	db, err := InitDB(dbPath)
	if err != nil {
		return err
	}
	return SaveGameLogs(db, gameStateID, messages)
}

// LoadGameLogs loads all saved log messages for a given gameStateID.
func LoadGameLogs(db *gorm.DB, gameStateID uint) ([]string, error) {
	if gameStateID == 0 {
		return nil, nil
	}
	var logs []GameLog
	if err := db.Where("game_state_id = ?", gameStateID).Order("id asc").Find(&logs).Error; err != nil {
		return nil, err
	}
	var res []string
	for _, l := range logs {
		res = append(res, l.Message)
	}
	return res, nil
}

// LoadGameLogsFromPath opens DB at dbPath and loads log messages.
func LoadGameLogsFromPath(dbPath string, gameStateID uint) ([]string, error) {
	if dbPath == "" {
		dbPath = DefaultDBPath()
	}
	db, err := InitDB(dbPath)
	if err != nil {
		return nil, err
	}
	return LoadGameLogs(db, gameStateID)
}

// ClearGameLogs deletes all saved log messages for the given gameStateID.
func ClearGameLogs(db *gorm.DB, gameStateID uint) error {
	if gameStateID == 0 {
		return nil
	}
	return db.Where("game_state_id = ?", gameStateID).Delete(&GameLog{}).Error
}

// ClearGameLogsToPath deletes all saved log messages in DB at dbPath.
func ClearGameLogsToPath(dbPath string, gameStateID uint) error {
	if dbPath == "" {
		dbPath = DefaultDBPath()
	}
	db, err := InitDB(dbPath)
	if err != nil {
		return err
	}
	return ClearGameLogs(db, gameStateID)
}

// SavePartyHeroes saves updated stats/HP for party heroes.
func SavePartyHeroes(db *gorm.DB, partyID uint, heroes []Hero) error {
	for i := range heroes {
		if heroes[i].PartyID == nil {
			heroes[i].PartyID = &partyID
		}
		if heroes[i].ID != 0 {
			_ = db.Save(&heroes[i]).Error
		}
	}
	return nil
}

// SavePartyHeroesToPath opens DB at dbPath and saves party hero stats.
func SavePartyHeroesToPath(dbPath string, partyID uint, heroes []Hero) error {
	if dbPath == "" {
		dbPath = DefaultDBPath()
	}
	db, err := InitDB(dbPath)
	if err != nil {
		return err
	}
	return SavePartyHeroes(db, partyID, heroes)
}


