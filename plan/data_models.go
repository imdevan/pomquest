package main

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Noun struct {
	gorm.Model
	Name         sql.String
	Age          sql.NullInt16 
	Hp	sql.sqlInt
}

// possibly an enum?
type Stat struct {
	Name 	sql.String
	Short	sql.String
	Value sql.Int16
}

// possibly an enum?
type DamageType struct {
	Name 	sql.String
	HitRoll sql.Int
	DmgRoll sql.Int
}

// possibly an enum?
type Class struct {
	Name 	sql.String
	Profficiency Stat[]
	DmgType DamageType
}
 
type Stats struct {
	Str sql.Int16 
	Dex sql.Int16
	Cha sql.Int16 
	Sur sql.Int16
	Con sql.Int16 
	Wil sql.Int16
	Int sql.Int16 
	Per sql.Int16
	Hel sql.Int16 
	Luc sql.Int16
}

type Horde struct {
	gorm.Model
	TreasureChests sql.Int16
	MagicItems sql.Int16
}

type Player struct {
	gorm.Model
	LoopsStarted 		sql.Int16
	LoopsCompleted 	sql.Int16
	TimeFocused			sql.Int16
}

type NPC struct {
	Noun
}

type Hero struct {
	Noun
	Stats
}

type Party scrut {
	potions  *uint
	coin  *uint
	parts  *uint
	parts  *uint
	
}

type NPC struct{}

type Monster struct {
	gorm.Model
}

type Quest struct{}

type Location struct {

}

type GameState struct {
	gorm.Model
	location  id

}

