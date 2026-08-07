package icon

// Mode defines the icon rendering variant (Nerd Font or Emoji).
type Mode int

const (
	ModeNerdFont Mode = iota
	ModeEmoji
)

// IconKey represents a symbolic icon name mapped to UTF / Emoji representations.
type IconKey string

const (
	KeyAttack     IconKey = "attack"
	KeyArmor      IconKey = "armor"
	KeyHP         IconKey = "hp"
	KeyDice       IconKey = "dice"
	KeyWeather    IconKey = "weather"
	KeyTerrain    IconKey = "terrain"
	KeyQuest      IconKey = "quest"
	KeyGold       IconKey = "gold"
	KeySilver     IconKey = "silver"
	KeyCopper     IconKey = "copper"
	KeyFood       IconKey = "food"
	KeyOx         IconKey = "ox"
	KeySpareParts IconKey = "spare_parts"
	KeyMagicItem  IconKey = "magic_item"
	KeyGem        IconKey = "gem"
	KeyChest      IconKey = "chest"
	KeySick       IconKey = "sick"
)

// IconPair stores both the Nerd Font UTF icon equivalent and the Emoji variant.
type IconPair struct {
	NerdFont string
	Emoji    string
}

var registry = map[IconKey]IconPair{
	KeyAttack:     {NerdFont: "󰓥", Emoji: "⚔️"},
	KeyArmor:      {NerdFont: "🛡", Emoji: "🛡️"},
	KeyHP:         {NerdFont: "", Emoji: "❤️"},
	KeyDice:       {NerdFont: "󰇍", Emoji: "🎲"},
	KeyWeather:    {NerdFont: "☀️", Emoji: "☀️"},
	KeyTerrain:    {NerdFont: "🌲", Emoji: "🌲"},
	KeyQuest:      {NerdFont: "🗡", Emoji: "🗡️"},
	KeyGold:       {NerdFont: "🪙", Emoji: "🪙"},
	KeySilver:     {NerdFont: "🥈", Emoji: "🥈"},
	KeyCopper:     {NerdFont: "🥉", Emoji: "🥉"},
	KeyFood:       {NerdFont: "🍖", Emoji: "🍖"},
	KeyOx:         {NerdFont: "🐂", Emoji: "🐂"},
	KeySpareParts: {NerdFont: "⚙️", Emoji: "⚙️"},
	KeyMagicItem:  {NerdFont: "✨", Emoji: "✨"},
	KeyGem:        {NerdFont: "💎", Emoji: "💎"},
	KeyChest:      {NerdFont: "📦", Emoji: "📦"},
	KeySick:       {NerdFont: "🤢", Emoji: "🤢"},
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
