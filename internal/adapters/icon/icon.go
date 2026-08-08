package icon

import "strings"

// Mode defines the icon rendering variant (Nerd Font or Emoji).
type Mode int

const (
	ModeNerdFont Mode = iota
	ModeEmoji
)

// IconKey represents a symbolic icon name mapped to UTF / Emoji representations.
type IconKey string

const (
	KeyAttack       IconKey = "attack"
	KeyArmor        IconKey = "armor"
	KeyHP           IconKey = "hp"
	KeyHeart10Full  IconKey = "heart_10_full"
	KeyHeart10Half  IconKey = "heart_10_half"
	KeyHeart10Empty IconKey = "heart_10_empty"
	KeyHeart5Full   IconKey = "heart_5_full"
	KeyHeart5Empty  IconKey = "heart_5_empty"
	KeyDice         IconKey = "dice"
	KeyWeather      IconKey = "weather"
	KeyTerrain      IconKey = "terrain"
	KeyQuest        IconKey = "quest"
	KeyGold         IconKey = "gold"
	KeySilver       IconKey = "silver"
	KeyCopper       IconKey = "copper"
	KeyFood         IconKey = "food"
	KeyOx           IconKey = "ox"
	KeySpareParts   IconKey = "spare_parts"
	KeyMagicItem    IconKey = "magic_item"
	KeyGem          IconKey = "gem"
	KeyChest        IconKey = "chest"
	KeySick         IconKey = "sick"
	KeyLoot         IconKey = "loot"
	KeyRope         IconKey = "rope"
	KeyBomb         IconKey = "bomb"
	KeyPotion       IconKey = "potion"
	KeyShovel       IconKey = "shovel"

	// Stat category keys
	KeyItemUsage IconKey = "item_usage"
	KeyDiet      IconKey = "diet"
	KeyGreed     IconKey = "greed"
	KeyPace      IconKey = "pace"

	// Class keys
	KeyClassFighter   IconKey = "class_fighter"
	KeyClassWizard    IconKey = "class_wizard"
	KeyClassWarlock   IconKey = "class_warlock"
	KeyClassRanger    IconKey = "class_ranger"
	KeyClassRogue     IconKey = "class_rogue"
	KeyClassCleric    IconKey = "class_cleric"
	KeyClassPaladin   IconKey = "class_paladin"
	KeyClassBard      IconKey = "class_bard"
	KeyClassDruid     IconKey = "class_druid"
	KeyClassBarbarian IconKey = "class_barbarian"
	KeyClassMonk      IconKey = "class_monk"
	KeyClassGuard     IconKey = "class_guard"

	// Item usage keys
	KeyItemUsageNone     IconKey = "item_usage_none"
	KeyItemUsageSparse   IconKey = "item_usage_sparse"
	KeyItemUsageBalanced IconKey = "item_usage_balanced"
	KeyItemUsageHeavy    IconKey = "item_usage_heavy"

	// Diet keys
	KeyDietNone     IconKey = "diet_none"
	KeyDietSparse   IconKey = "diet_sparse"
	KeyDietBalanced IconKey = "diet_balanced"
	KeyDietFilling  IconKey = "diet_filling"

	// Greed keys
	KeyGreedNone     IconKey = "greed_none"
	KeyGreedSparse   IconKey = "greed_sparse"
	KeyGreedBalanced IconKey = "greed_balanced"
	KeyGreedHeavy    IconKey = "greed_heavy"

	// Pace keys
	KeyPaceSlow     IconKey = "pace_slow"
	KeyPaceBalanced IconKey = "pace_balanced"
	KeyPaceFast     IconKey = "pace_fast"
)

// IconPair stores both the Nerd Font UTF icon equivalent and the Emoji variant.
type IconPair struct {
	NerdFont string
	Emoji    string
}

var registry = map[IconKey]IconPair{
	KeyAttack:       {NerdFont: "󰓥", Emoji: "⚔️"},
	KeyArmor:        {NerdFont: "🛡", Emoji: "🛡️"},
	KeyHP:           {NerdFont: "", Emoji: "❤️"},
	KeyHeart10Full:  {NerdFont: "󰋑", Emoji: "❤️"},
	KeyHeart10Half:  {NerdFont: "󰛞", Emoji: "💔"},
	KeyHeart10Empty: {NerdFont: "󰋕", Emoji: "🖤"},
	KeyHeart5Full:   {NerdFont: "󰛟", Emoji: "💗"},
	KeyHeart5Empty:  {NerdFont: "󰛠", Emoji: "🤍"},
	KeyDice:         {NerdFont: "󰇍", Emoji: "🎲"},
	KeyWeather:      {NerdFont: "☀️", Emoji: "☀️"},
	KeyTerrain:      {NerdFont: "🌲", Emoji: "🌲"},
	KeyQuest:        {NerdFont: "🗡", Emoji: "🗡️"},
	KeyMagicItem:    {NerdFont: "✨", Emoji: "✨"},
	KeyGem:          {NerdFont: "💎", Emoji: "💎"},
	KeyChest:        {NerdFont: "📦", Emoji: "📦"},
	KeySick:         {NerdFont: "🤢", Emoji: "🤢"},

	KeyGold:       {NerdFont: "🪙", Emoji: "🪙"},
	KeyCopper:     {NerdFont: "🥉", Emoji: "🥉"},
	KeySilver:     {NerdFont: "🥈", Emoji: "🥈"},
	KeyFood:       {NerdFont: "󱐟", Emoji: "🍖"},
	KeyLoot:       {NerdFont: "󱢏", Emoji: "💰"},
	KeySpareParts: {NerdFont: "", Emoji: "⚙️"},
	KeyOx:         {NerdFont: "", Emoji: "🐂"},

	KeyRope:   {NerdFont: "󰼃", Emoji: "🪢"},
	KeyBomb:   {NerdFont: "󰚑", Emoji: "💣"},
	KeyPotion: {NerdFont: "", Emoji: "🧪"},
	KeyShovel: {NerdFont: "󰜐", Emoji: "⛏️"},

	KeyClassFighter:   {NerdFont: "⚔", Emoji: "⚔️"},
	KeyClassWizard:    {NerdFont: "", Emoji: "🧙"},
	KeyClassWarlock:   {NerdFont: "󰱯", Emoji: "👁️"},
	KeyClassRanger:    {NerdFont: "󱡁", Emoji: "🏹"},
	KeyClassRogue:     {NerdFont: "🗡", Emoji: "🗡️"},
	KeyClassCleric:    {NerdFont: "⚕", Emoji: "⚕️"},
	KeyClassPaladin:   {NerdFont: "󰳊", Emoji: "🛡️"},
	KeyClassBard:      {NerdFont: "󰝱", Emoji: "🪕"},
	KeyClassDruid:     {NerdFont: "", Emoji: "🌿"},
	KeyClassBarbarian: {NerdFont: "󰿖", Emoji: "💪"},
	KeyClassMonk:      {NerdFont: "", Emoji: "✊"},
	KeyClassGuard:     {NerdFont: "󱢾", Emoji: "🛡️"},

	KeyItemUsageNone:     {NerdFont: "", Emoji: "🚫"},
	KeyItemUsageSparse:   {NerdFont: "", Emoji: "💧"},
	KeyItemUsageBalanced: {NerdFont: "", Emoji: "⚖️"},
	KeyItemUsageHeavy:    {NerdFont: "", Emoji: "🔥"},

	KeyDietNone:     {NerdFont: "", Emoji: "🚫"},
	KeyDietSparse:   {NerdFont: "", Emoji: "🍃"},
	KeyDietBalanced: {NerdFont: "", Emoji: "🍲"},
	KeyDietFilling:  {NerdFont: "", Emoji: "🍖"},

	KeyGreedNone:     {NerdFont: "", Emoji: "🕊️"},
	KeyGreedSparse:   {NerdFont: "", Emoji: "🪙"},
	KeyGreedBalanced: {NerdFont: "", Emoji: "💰"},
	KeyGreedHeavy:    {NerdFont: "", Emoji: "🤑"},

	KeyPaceSlow:     {NerdFont: "", Emoji: "🐢"},
	KeyPaceBalanced: {NerdFont: "󰜎", Emoji: "🚶"},
	KeyPaceFast:     {NerdFont: "󰑮", Emoji: "🏃"},

	KeyItemUsage: {NerdFont: "", Emoji: "🎒"},
	KeyDiet:      {NerdFont: "", Emoji: "🍽️"},
	KeyPace:      {NerdFont: "󰖃", Emoji: "⏱️"},
	KeyGreed:     {NerdFont: "", Emoji: "👑"},
}

// Get returns the icon string for a given key and display mode.
func Get(key IconKey, mode Mode) string {
	pair, ok := registry[key]
	if !ok {
		return string(key)
	}
	if mode == ModeEmoji && pair.Emoji != "" {
		return pair.Emoji
	}
	return pair.NerdFont
}

// GetNerdFont returns the Nerd Font UTF icon equivalent.
func GetNerdFont(key IconKey) string {
	return Get(key, ModeNerdFont)
}

// GetEmoji returns the Emoji icon variant.
func GetEmoji(key IconKey) string {
	return Get(key, ModeEmoji)
}

// GetItemUsageIcon returns the icon string for an item usage value.
func GetItemUsageIcon(val string, mode Mode) string {
	return Get(GetItemUsageIconKey(val), mode)
}

func GetItemUsageIconKey(val string) IconKey {
	switch val {
	case "none":
		return KeyItemUsageNone
	case "sparse":
		return KeyItemUsageSparse
	case "balanced":
		return KeyItemUsageBalanced
	case "heavy":
		return KeyItemUsageHeavy
	default:
		return KeyItemUsageNone
	}
}

// GetDietIcon returns the icon string for a diet value.
func GetDietIcon(val string, mode Mode) string {
	return Get(GetDietIconKey(val), mode)
}

func GetDietIconKey(val string) IconKey {
	switch val {
	case "none":
		return KeyDietNone
	case "sparse":
		return KeyDietSparse
	case "balanced":
		return KeyDietBalanced
	case "filling":
		return KeyDietFilling
	default:
		return KeyDietBalanced
	}
}

// GetGreedIcon returns the icon string for a greed value.
func GetGreedIcon(val string, mode Mode) string {
	return Get(GetGreedIconKey(val), mode)
}

func GetGreedIconKey(val string) IconKey {
	switch val {
	case "none":
		return KeyGreedNone
	case "sparse", "low":
		return KeyGreedSparse
	case "balanced":
		return KeyGreedBalanced
	case "heavy":
		return KeyGreedHeavy
	default:
		return KeyGreedNone
	}
}

// GetPaceIcon returns the icon string for a pace value.
func GetPaceIcon(val string, mode Mode) string {
	return Get(GetPaceIconKey(val), mode)
}

func GetPaceIconKey(val string) IconKey {
	switch val {
	case "slow":
		return KeyPaceSlow
	case "balanced":
		return KeyPaceBalanced
	case "fast":
		return KeyPaceFast
	default:
		return KeyPaceFast
	}
}

// GetClassIcon returns the icon string for a hero class name or key.
func GetClassIcon(val string, mode Mode) string {
	return Get(GetClassIconKey(val), mode)
}

func GetClassIconKey(val string) IconKey {
	switch strings.ToLower(strings.TrimSpace(val)) {
	case "fighter", "warrior", "knight", "class_fighter":
		return KeyClassFighter
	case "wizard", "mage", "sorcerer", "class_wizard":
		return KeyClassWizard
	case "warlock", "class_warlock":
		return KeyClassWarlock
	case "ranger", "hunter", "archer", "class_ranger":
		return KeyClassRanger
	case "rogue", "rouge", "thief", "assassin", "class_rogue":
		return KeyClassRogue
	case "cleric", "priest", "healer", "class_cleric":
		return KeyClassCleric
	case "paladin", "class_paladin":
		return KeyClassPaladin
	case "bard", "class_bard":
		return KeyClassBard
	case "druid", "class_druid":
		return KeyClassDruid
	case "barbarian", "class_barbarian":
		return KeyClassBarbarian
	case "monk", "class_monk":
		return KeyClassMonk
	case "guard", "tank", "defender", "class_guard":
		return KeyClassGuard
	default:
		return IconKey(val)
	}
}

// Legacy editor & file icons:
type Icon string

const (
	Bookmarks Icon = "\uf097"
	Tmux      Icon = "\uebc8"
	Nvim      Icon = "\uf36f"
	Vim       Icon = "\uf36f"
	File      Icon = "\uf15b"
	Script    Icon = "\ue691"
	Shell     Icon = "\uf489"
)

func (i Icon) String() string {
	return string(i)
}

func GetEditorIcon(editor string) Icon {
	switch editor {
	case "nvim", "neovim":
		return Nvim
	case "vim":
		return Vim
	default:
		return ""
	}
}
