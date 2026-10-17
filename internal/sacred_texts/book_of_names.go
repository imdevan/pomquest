package sacred_texts

// characterNames maps race -> gender -> name data.
var characterNames = map[Race]map[Gender]nameData{
	RaceHuman: {
		GenderMale: {
			Prefixes: []string{
				"Gar", "Bran", "Ald", "Ced", "Thor", "Ed", "Wil", "Ric", "Arn", "God",
				"Ead", "Alf", "Har", "Leon", "Os", "Rod", "Sig", "Ulf", "Wulf", "Bald",
			},
			Suffixes: []string{
				"ric", "don", "win", "mund", "ard", "bert", "fred", "wald", "olf", "red",
				"gar", "stan", "helm", "mar", "rik", "wen", "ard", "mund", "orn", "var",
			},
		},
		GenderFemale: {
			Prefixes: []string{
				"Ari", "Ela", "Bri", "Gwen", "Isa", "Ros", "Miri", "Alyn", "Hed", "Wyn",
				"Ath", "Bere", "Clo", "Edith", "Hild", "Leof", "Maud", "Sif", "Tild", "Yse",
			},
			Suffixes: []string{
				"ella", "wen", "aine", "dra", "lyn", "isa", "ora", "iel", "eth", "ild",
				"ina", "ith", "ola", "une", "ira", "ene", "ary", "ess", "iel", "ane",
			},
		},
	},
	RaceElf: {
		GenderMale: {
			Prefixes: []string{
				"Ael", "Thal", "Lor", "Gal", "Fen", "Cel", "Syl", "Ael", "Cael", "Tir",
				"Ear", "Fin", "Gil", "Hal", "Leg", "Mal", "Nim", "Orl", "Riv", "Val",
			},
			Suffixes: []string{
				"aril", "ion", "dor", "wen", "ith", "ael", "orn", "as", "iel", "an",
				"dir", "las", "mir", "nor", "oth", "ryn", "thar", "uin", "wen", "yr",
			},
		},
		GenderFemale: {
			Prefixes: []string{
				"Ar", "El", "Luth", "Nym", "Syl", "Aer", "Ith", "Lia", "Gael", "Mir",
				"Cel", "Eil", "Fae", "Ier", "Lael", "Nar", "Quel", "Tae", "Vael", "Yor",
			},
			Suffixes: []string{
				"wen", "riel", "ara", "iel", "ith", "anwe", "ora", "ael", "ala", "iel",
				"ari", "dis", "eth", "iel", "lyn", "mir", "nai", "sel", "thi", "ura",
			},
		},
	},
	RaceDwarf: {
		GenderMale: {
			Prefixes: []string{
				"Dur", "Thor", "Brom", "Gim", "Bal", "Oin", "Dwal", "Gror", "Nor", "Bol",
				"Bar", "Dor", "Fal", "Glor", "Kor", "Mur", "Nar", "Skog", "Thr", "Vor",
			},
			Suffixes: []string{
				"in", "li", "ik", "ur", "ok", "rin", "ak", "im", "or", "un",
				"din", "grim", "mund", "rak", "dur", "bek", "gar", "dak", "mir", "tun",
			},
		},
		GenderFemale: {
			Prefixes: []string{
				"Dis", "Hel", "Bry", "Fre", "Gur", "Ket", "Myr", "Thr", "Aud", "Skj",
				"Ast", "Dag", "Eid", "Gerd", "Hild", "Ing", "Kol", "Ragna", "Sig", "Tor",
			},
			Suffixes: []string{
				"dis", "hild", "ia", "lin", "ra", "ya", "nis", "ila", "rid", "va",
				"bra", "dur", "gar", "ka", "mir", "nir", "run", "tha", "und", "vig",
			},
		},
	},
	RaceHobbit: {
		GenderMale: {
			Prefixes: []string{
				"Bil", "Fro", "Sam", "Pip", "Mer", "Dro", "Tol", "Ham", "Per", "Fol",
				"Ban", "Bun", "Dim", "Fal", "Ger", "Hob", "Jol", "Mos", "Nor", "Wil",
			},
			Suffixes: []string{
				"bo", "do", "go", "ry", "in", "fast", "wise", "ard", "co", "lo",
				"ric", "by", "ton", "kins", "well", "dale", "ford", "wick", "den", "mow",
			},
		},
		GenderFemale: {
			Prefixes: []string{
				"Ros", "Prim", "Lob", "Cam", "Eri", "Dais", "Pop", "Ama", "Bel", "Jas",
				"Ang", "Clem", "Est", "Gar", "Hol", "Ivy", "Lil", "Myr", "Pan", "Tan",
			},
			Suffixes: []string{
				"ie", "ula", "elia", "a", "ina", "ey", "ola", "ily", "ia", "ina",
				"ry", "sy", "da", "belle", "ine", "ony", "wyn", "lee", "beth", "rose",
			},
		},
	},
	RaceGoblin: {
		GenderMale: {
			Prefixes: []string{
				"Griz", "Snag", "Mog", "Nak", "Gob", "Skrit", "Zag", "Rik", "Nib", "Grot",
				"Blat", "Crik", "Dreg", "Flik", "Gnar", "Jix", "Krik", "Pug", "Spit", "Wrek",
			},
			Suffixes: []string{
				"nak", "gob", "ik", "ot", "uk", "az", "ek", "it", "ok", "ug",
				"ax", "bit", "fiz", "gak", "nit", "pok", "riz", "snik", "tik", "zak",
			},
		},
		GenderFemale: {
			Prefixes: []string{
				"Nix", "Shiv", "Grit", "Mez", "Snik", "Pik", "Raz", "Tik", "Zil", "Kri",
				"Blix", "Driz", "Fizz", "Glim", "Hix", "Jez", "Lix", "Rix", "Skiz", "Vex",
			},
			Suffixes: []string{
				"ka", "sha", "zi", "na", "ra", "li", "ta", "ki", "la", "si",
				"ax", "bi", "da", "fi", "ga", "ni", "pi", "ri", "va", "za",
			},
		},
	},
	RaceOrc: {
		GenderMale: {
			Prefixes: []string{
				"Grom", "Urg", "Mog", "Gor", "Krag", "Zug", "Throk", "Bol", "Drog", "Skar",
				"Azog", "Brug", "Durg", "Gash", "Gruk", "Lug", "Murg", "Narg", "Shag", "Vrog",
			},
			Suffixes: []string{
				"mash", "gul", "gor", "rok", "dak", "thar", "nak", "ug", "ash", "rak",
				"bur", "dag", "gash", "kul", "lug", "mok", "rag", "thak", "urk", "zog",
			},
		},
		GenderFemale: {
			Prefixes: []string{
				"Gra", "Ura", "Shel", "Mog", "Gor", "Bor", "Lur", "Yag", "Tha", "Rug",
				"Aga", "Bul", "Dur", "Gul", "Hag", "Kra", "Muz", "Nar", "Sha", "Vra",
			},
			Suffixes: []string{
				"gra", "sha", "bra", "ka", "za", "ra", "da", "na", "la", "ma",
				"ga", "tha", "ba", "dra", "gha", "kha", "mba", "rza", "sha", "zha",
			},
		},
	},
}

// BookOfLocationNames contains a curated collection of location base names.
var BookOfLocationNames = []string{
	"Izelliador",
	"Malakor",
	"Zalthor",
	"Gloomhaven",
	"Vaelgor",
	"Kardoth",
	"Thulgor",
	"Baelroth",
	"Dreadfang",
	"Shadowveil",
	"Obsidian",
	"Bloodspire",
	"Voragath",
	"Kaelmon",
	"Drakfell",
	"Arakmire",
	"Xalnoth",
	"Grimhold",
	"Frostcrag",
	"Ironreach",
	"Stormforge",
	"Voidvault",
	"Nightmar",
	"Skymor",
	"Emberfall",
	"Blightwood",
	"Ashridge",
	"Morvath",
}

// LocationPrefixes allows generating donjon-style location names.
var LocationPrefixes = []string{
	"Bael", "Zal", "Kard", "Thul", "Vael", "Gloom", "Dread", "Gor", "Mael", "Mor",
	"Rath", "Drak", "Khar", "Grim", "Vor", "Shadow", "Ash", "Storm", "Iron", "Frost",
	"Void", "Night", "Mal", "Izel", "Arak", "Xal", "Kael", "Bal", "Krag", "Vex",
}

// LocationSuffixes allows generating donjon-style location names.
var LocationSuffixes = []string{
	"roth", "thor", "doth", "gor", "haven", "fang", "veil", "spire", "mire", "hold",
	"keep", "fell", "crag", "reach", "depth", "forge", "gate", "ridge", "vault", "mar",
	"dor", "mon", "gath", "noth", "liador", "heim", "crest", "glen",
}
