package app

import (
	"time"

	"gorm.io/gorm"
)

// Stats represents character attribute scores used for rolls and scaling.
type Stats struct {
	Str int `json:"str" gorm:"default:0"`
	Dex int `json:"dex" gorm:"default:0"`
	Cha int `json:"cha" gorm:"default:0"`
	Sur int `json:"sur" gorm:"default:0"`
	Con int `json:"con" gorm:"default:0"`
	Wil int `json:"wil" gorm:"default:0"`
	Int int `json:"int" gorm:"default:0"`
	Per int `json:"per" gorm:"default:0"`
	Hel int `json:"hel" gorm:"default:0"`
	Luc int `json:"luc" gorm:"default:0"`
}

// Noun represents any named entity with a visual icon/sprite and database primary key.
type Noun struct {
	gorm.Model
	Name string `gorm:"size:64;not null"`
	Icon string `gorm:"size:16"`
}

// Creature represents an entity with combat stats (HP, Attack, Armor).
type Creature struct {
	Noun
	HP     int `gorm:"default:10"`
	MaxHP  int `gorm:"default:10"`
	Attack int `gorm:"default:1"`
	Armor  int `gorm:"default:0"`
}

// Hero represents an individual party member / adventurer.
type Hero struct {
	Creature
	PartyID *uint  `gorm:"index"`
	Class   string `gorm:"size:32;not null"`
	Level   int    `gorm:"default:1"`
	XP      int    `gorm:"default:0"`
	Perk    string `gorm:"size:64"`
	Status  string `gorm:"size:64;default:'Healthy'"` // Healthy, Tired, Sick, Dysentery, etc.
	IsDead  bool   `gorm:"default:false"`
	Stats   `gorm:"embedded;embeddedPrefix:stat_"`
}

// ClassDef represents predefined RPG class archetypes and base stats.
type ClassDef struct {
	gorm.Model
	Name            string `gorm:"size:32;uniqueIndex;not null"`
	Description     string `gorm:"size:255"`
	PrimaryStat     string `gorm:"size:16"`
	PerkName        string `gorm:"size:64"`
	PerkDescription string `gorm:"size:255"`
	Stats           `gorm:"embedded;embeddedPrefix:base_stat_"`
}

// Monster represents an enemy encountered on quests.
type Monster struct {
	Creature
	TargetPrio string `gorm:"size:32;default:'tank'"` // tank, int, dex, lowest_hp
	XPValue    int    `gorm:"default:10"`
	LootGold   int    `gorm:"default:1"`
	Biome      string `gorm:"size:64"`
}

// NPC represents non-player characters in towns or along quest paths.
type NPC struct {
	Noun
	Role     string `gorm:"size:32"` // Merchant, Innkeeper, QuestGiver, etc.
	Dialogue string `gorm:"type:text"`
}

// Party represents the active group of heroes and shared supplies.
type Party struct {
	gorm.Model
	Name       string `gorm:"size:64;default:'Adventuring Party'"`
	Gold       int    `gorm:"default:0"`
	Silver     int    `gorm:"default:0"`
	Copper     int    `gorm:"default:0"`
	FoodLbs    int    `gorm:"default:100"`
	Diet       string `gorm:"size:32;default:'balanced'"` // sparse, balanced, heavy
	Pace       string `gorm:"size:32;default:'fast'"`     // slow, balanced, fast
	SpareParts int    `gorm:"default:2"`
	Oxen       int    `gorm:"default:2"`
	Regents    int    `gorm:"default:0"`
	Potions    int    `gorm:"default:0"`
	Poisons    int    `gorm:"default:0"`
	Loot       int    `gorm:"default:0"`
	Rope       int    `gorm:"default:0"`
	Bombs      int    `gorm:"default:0"`
	Shovels    int    `gorm:"default:0"`
	ItemUsage  string `gorm:"size:32;default:'none'"` // none, sparse, balanced, heavy
	Greed      string `gorm:"size:32;default:'none'"` // none, low, med, high
	Heroes     []Hero `gorm:"foreignKey:PartyID"`
}

// Hoard represents player's overall accumulated hoard / score across loops.
type Hoard struct {
	gorm.Model
	TreasureChests int `gorm:"default:0"`
	MagicItems     int `gorm:"default:0"`
	GoldTotal      int `gorm:"default:0"`
}

// Quest represents an adventure undertaking.
type Quest struct {
	gorm.Model
	Title           string `gorm:"size:128;not null"`
	Description     string `gorm:"type:text"`
	Biome           string `gorm:"size:64"`
	PathLength      int    `gorm:"default:100"`
	CurrentDistance int    `gorm:"default:0"`
	Weather         string `gorm:"size:32;default:'Sunny'"`
	Status          string `gorm:"size:32;default:'Available'"` // Available, Active, Completed, Failed
	Difficulty      int    `gorm:"default:1"`
	RewardGold      int    `gorm:"default:0"`
}

// Location represents a biome or destination.
type Location struct {
	gorm.Model
	Name               string  `gorm:"size:64;not null"`
	Biome              string  `gorm:"size:64"`
	Description        string  `gorm:"type:text"`
	DifficultyModifier float64 `gorm:"default:1.0"`
}

// Encounter represents random events or combat occurrences on a quest.
type Encounter struct {
	gorm.Model
	QuestID     uint   `gorm:"index"`
	Type        string `gorm:"size:32"` // Combat, Trap, Rest, Merchant, Event
	Description string `gorm:"type:text"`
	Outcome     string `gorm:"type:text"`
}

// GameState holds the active session state linking quest, party, Hoard, and pomodoro tracking.
type GameState struct {
	gorm.Model
	CurrentQuestID       *uint     `gorm:"index"`
	CurrentQuest         *Quest    `gorm:"foreignKey:CurrentQuestID"`
	PartyID              *uint     `gorm:"index"`
	Party                *Party    `gorm:"foreignKey:PartyID"`
	HoardID              *uint     `gorm:"index"`
	Hoard                *Hoard    `gorm:"foreignKey:HoardID"`
	CurrentLocationID    *uint     `gorm:"index"`
	CurrentLocation      *Location `gorm:"foreignKey:CurrentLocationID"`
	LoopsStarted         int       `gorm:"default:0"`
	LoopsCompleted       int       `gorm:"default:0"`
	TimeFocusedMinutes   int       `gorm:"default:0"`
	CurrentPomodoroState string    `gorm:"size:32;default:'idle'"` // working, short_break, long_break, idle
	IsActive             bool      `gorm:"default:true"`
}

// Player represents cumulative global player stats across all games.
type Player struct {
	gorm.Model
	Name            string `gorm:"size:64;default:'Dark Lord'"`
	TotalRuns       int    `gorm:"default:0"`
	QuestsCompleted int    `gorm:"default:0"`
	HeroesRecruited int    `gorm:"default:0"`
	HeroesLost      int    `gorm:"default:0"`
	LoopsStarted    int    `gorm:"default:0"`
	LoopsCompleted  int    `gorm:"default:0"`
	TimeFocused     int    `gorm:"default:0"`
}

// GameLog tracks event history for UI scroll and persistence.
type GameLog struct {
	gorm.Model
	GameStateID uint      `gorm:"index"`
	Message     string    `gorm:"type:text;not null"`
	Category    string    `gorm:"size:32;default:'general'"` // combat, item, quest, pomodoro
	Timestamp   time.Time `gorm:"autoCreateTime"`
}
