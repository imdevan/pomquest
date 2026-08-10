package app

import "github.com/pomquest/internal/adapters/icon"

// ItemDef represents a purchasable item in the General Store book of items.
type ItemDef struct {
	ID             string
	Name           string
	Unit           string
	BaseCostCopper int
	DisplayCost    string
	Description    string
	IconKey        icon.IconKey
}

// BookOfItems returns the list of all items available at Soand Sow's General Store.
func BookOfItems() []ItemDef {
	return []ItemDef{
		{
			ID:             "food",
			Name:           "Food",
			Unit:           "1 lb",
			BaseCostCopper: 10, // 10c / lb
			DisplayCost:    "10c / lb",
			Description:    "Basic provisions to sustain party members on journeys.",
			IconKey:        icon.KeyFood,
		},
		{
			ID:             "spare_parts",
			Name:           "Spare Parts",
			Unit:           "1 part",
			BaseCostCopper: 60, // 60c / part
			DisplayCost:    "60c / part",
			Description:    "Essential wagon parts to repair breakdowns on the trail.",
			IconKey:        icon.KeySpareParts,
		},
		{
			ID:             "ox",
			Name:           "Ox Yoke (3 Oxen)",
			Unit:           "1 yoke",
			BaseCostCopper: 200, // 20s = 200c per yoke (3 ox)
			DisplayCost:    "20s / yoke (3 ox)",
			Description:    "Strong oxen to pull wagons through rugged biomes.",
			IconKey:        icon.KeyOx,
		},
		{
			ID:             "potions",
			Name:           "Healing Potions",
			Unit:           "1 pot",
			BaseCostCopper: 200, // 20s = 200c / pot
			DisplayCost:    "20s / pot",
			Description:    "Restores HP and cures status ailments during quests.",
			IconKey:        icon.KeyPotion,
		},
		{
			ID:             "bombs",
			Name:           "Bombs",
			Unit:           "1 bomb",
			BaseCostCopper: 50, // 50c / bomb
			DisplayCost:    "50c / bomb",
			Description:    "Explosives to clear obstacles or blast groups of enemies.",
			IconKey:        icon.KeyBomb,
		},
		{
			ID:             "shovels",
			Name:           "Shovels",
			Unit:           "1 shovel",
			BaseCostCopper: 20, // 20c / shovel
			DisplayCost:    "20c / shovel",
			Description:    "Dig through rockslides and unearth hidden treasure.",
			IconKey:        icon.KeyShovel,
		},
		{
			ID:             "rope",
			Name:           "Rope",
			Unit:           "50ft",
			BaseCostCopper: 10, // 10c / 50ft
			DisplayCost:    "10c / 50ft",
			Description:    "Heavy hemp rope to scale cliffs and cross chasms.",
			IconKey:        icon.KeyRope,
		},
	}
}
