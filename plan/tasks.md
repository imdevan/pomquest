# Context

## Goal
The goal of this projects to create tui project oult

Concept layout more thoroughly in plan/mission.md

# Features
## Feature 1: define data models
- [x] 1.1 define / refine gorm data models for game in internal/app/data_models.go.
- do not integrate into app yet. 
- data models will be used to represent everything from game state, to monsters, to heroes

  

## Feature 2: define sqlite db architecture
- [x] 2.1 document findings in plan/sqlite-discovery.md
- [x] 2.2 analyze what should be stored vs embedded enums/code templates
- [x] 2.3 define simple migration patterns (GORM AutoMigrate & pragmas)
- [x] 2.4 define how app interacts with sqlite db (GameRepository layer & max open connections)
- [x] 2.5 define handling for early session termination (WAL mode, micro-saves, recovery prompt)
- [x] 2.6 outline additional architectural questions (multi-instance locking, log pruning, test mocks)

## Feature 3: define ui
- [x] 3.1 create map of user flows based on plan/screens.md, plan/screens_simple.md, and plan/initial-path.md in plan/define-ui-discovery.md
- [x] 3.2 tree diagram of ui architecture showing generic lib components vs domain specific components for shared theming in plan/define-ui-discovery.md

## Feature 4: simple dashboard screen
- [x] 4.1 implement hearth dashboard screen TUI component in internal/ui/dashboard.go
- [x] 4.2 wire root CLI command in cmd/pomquest/root.go to trigger dashboard screen on execution
- [x] 4.3 add unit tests in internal/ui/dashboard_test.go verifying navigation and option selection
## Feature 5: simple questing screen
- [x] 5.1 create generic UI components for progress bar, log view, hero card, and party supplies bar in internal/ui/
- [x] 5.2 implement questing screen model and prototype components with static dummy data in internal/screens/questing/
- [x] 5.3 wire ChoiceQuest transition from dashboard screen to questing screen in cmd/pomquest/root.go
- [x] 5.4 add unit tests in internal/screens/questing/questing_test.go verifying timer ticks and keypress events
- [x] 5.5 hero card fixes
  - [x] 5.5.1 hero cards should have 2 space in between to accoun for any card becoming active, and thus offset. 
  - [x] 5.5.2 hero cards should have padding of (1,2)
  - [x] 5.5.3 use 󰓥 for all hero card atk icons
  - [x] 5.5.4 atk / def alignment should be icon,space,number,space,icon,space,number
  - [x] 5.5.5 fix hero card content alignment (name & class top left, hero level top right, atk & def bottom left, hp bottom right)
  - [x] 5.5.6 active hero card - use two overlaping offset lipgloss boxes, rather than trying to manually redraw.
- [x] 5.6 party card fixes
  - [x] 5.6.1 partybar should be renamed to party card
  - [x] 5.6.2 party card should be composed of a lipgloss table component with 3 columns and one header row. 3d column should be right aligned. 

  |party | | |
  |---|---|---|
  |coins: 0g10s80c|diet: balanced| |
  |food: 2,000lbs|pace: fast| magic items: 0|
  |ox: 4| | gems: 0|
  |spare parts: 3| | tresure chests: 1|

## Feature 6: path screen
- [x] 6.1 create path screen
  - [x] 6.1.1 when a user begins a quest they should be directed to the path screen before the quest sceen
  - [x] 6.1.2 should still check for party size before starting
  - [x] 6.1.3 shows 3 generated dungeon names
  - [x] 6.1.4 use card format, similar to hero cards with box hover effect for selected item. 
- [x] 6.2 upon selection a Location should be generated and stored in the game states CurrentLocation

## Feature 7: connect quest screen
- [ ] 7.1 quest screen should show data  from sqlite db likely already in memory
- [ ] 7.2 hero cards should render based on current party heros
- [ ] 7.3 party card should show data form sqlite db
- [ ] 7.4 party card hearts should be based on current party hero health


## Feature 8: party card
- [x] 8.1 add team hearts ui component
  - [x] 8.1.1 takes max and current hp, renders hearts evaluated in 10 HP and 5 HP increments using icon.KeyHeart...
  - [x] 8.1.2 add team hearts to party texture card top-right

## Feature 9: create custom internal/ui/texturedCard
- [x] 9.1 create custom internal/ui/texturedCard component supporting embedded text on any of the 4 border corners with padding, margin, and styling properties
- [x] 9.2 use TexturedCard for party card container wrapping party table



## Feature 10: party model update
- [x] 10.1 add these to the party model:
  - [x] 10.1.1 small non-negative ints (loot, rope, bombs, potions, tshovels)
  - [x] 10.1.2 enum / string with fixed set of valid options (item usage: "none", "sparing", "balanced", "generous"; greed: "none", "low", "med", "high", "extreme")
- [x] 10.2 add those properties to the partycard
- [x] 10.3 add current party size and max party size
  - [x] 10.3.1 current party size defaults to 0, max party size defaults to 1 - both cached in sql db
- [ ] 10.4 add party name

## Feature 11: create icon package
- [x] 11.1 map strings to nerdfont utf equivalents
- [x] 11.2 some icons can have a nerd font and emoji variant
- [x] 11.3 create library of names mapped to icons
- [x] 11.4 replace all icons on quest screen components with icon package icons

## Feature 12: create icon library

## Feature 13: define class base stats
only warrior
str 2, con 2⌋
dmg: 1d6/2+1 => 0.5, 2, 2.5, 3, 3.5, 4
dmg: 1d6 % 4 = 2dmg! 
dmg: (1d6+2)/2 = 1.5, 2, 2.5, 3, 3.5, 4
dmg: ⌊(1d6+2)/2⌋ = 1, 2, 2, 3, 3, 4

maybe this
dmg: hits for 1d2 󰇋 (2) dmg!
## Feature 14: create randomization library
2d6 based system 2.6% crit chance - epic fail


## Feature 15: define quest, biome, path, weather tables
start with 1 quest, 1 biome, 1 path, 1 sunny weather


## Feature 16: define monster stat table

... maybe sql isn't needed for stats ext? unless, they are downloaded with first install... sqlite still good for saving player state. But does it make sense to ship all as embeded sqlite db?

... all stats could really be enums. what needs to saved is history, and player state... that could be file... sqlite would be good for large history storage, viewing previous runs, locations impacting each other, or revisited? 
- especially if using something like faker for names.


start with 1 monster, bats... 1hp 0def 1 atk

## Feature 17: add icons
- [x] 17.1 add heart icons in internal/adapters/icon/icon.go:74    
  - [x] 17.1.1 10hp heart 󰋑 
  - [x] 17.1.2 10hp - 5 heart 󰛞 
  - [x] 17.1.3 10hp - 10 heart 󰋕
  - [x] 17.1.4 5hp heart 󰛟
  - [x] 17.1.5 5hp empty heart 󰛠


## Feature 18: quest footer
- [x] 18.1 use progress bar for timer progress.
  - [x] 18.1.1 integrate github.com/charmbracelet/bubbles/progress model into ProgressBarComponent in internal/ui/progress.go
  - [x] 18.1.2 add unit test coverage in internal/ui/progress_test.go
- [x] 18.2 quest footer should be the full width of the questing container
- [x] 18.3 progress bar should be the full width of the container
- [x] 18.4 time should be aligned fully to right of container and not wrap line

## Feature 19: quest log
- [x] 19.1 replace quest log border with textureCard
- [x] 19.2 show quest log date top-left text on quest log texture card

## Feature 20: container
- [x] 20.1 create container that wraps quest screen components
  - [x] 20.1.1 all quest component sections flex to fit container width
  - [x] 20.1.2 set the width of the container to 96

## Feature 21: responsiveness
- [x] 21.1 adjust container and nested children to use internal/ui/responsive.go
  - [x] 21.1.1 on md: remove name from party status labels such that party card has `icon: value`
  - [x] 21.1.2 on md: remove space between icon atk and icon def on hero cards
  - [x] 21.1.3 on lg +: keep full label names and standard icon spacing
  - [x] 21.1.4 on sm -: show a message saying this terminal width is not currently supported with full pane border
    - [x] 19.1.5.1 update `QuestContainerComponent.Render` in `internal/screens/questing/container.go` to check for `sm -` breakpoints (`IsAtMost(ui.BreakpointSM)`)
    - [x] 19.1.5.2 render unsupported terminal width message with full pane border
    - [x] 19.1.5.3 add unit test `TestQuestContainerComponent_SmallScreenUnsupported` in `internal/screens/questing/container_test.go`

## Feature 22: fix screen transitions
- [x] 22.1 currently there is a blink where all text seems to be shifted to left side col when going from home to questing.
  - [x] 22.1.1 pass terminal dimensions across screen transitions to prevent uncentered rendering on frame 1
  - [x] 22.1.2 clear terminal screen on model Init to prevent transition flash


## Feature 23: add hoard screen
- [x] 23.1 hoard screen implementation
  - [x] 23.1.1 top section shows `treasure_hoard` ASCII art from `plan/ascii-art.md:575`
  - [x] 23.1.2 bottom section rendered in `TexturedCard` with `"Hoard"` header on top-left
  - [x] 23.1.3 3-column table showing total gold, total loot, total treasure, heroes hired, heroes lost, places traveled, quests started, quests completed, and time focused
  - [x] 23.1.4 define icon keys for hoard stats reusing gold, loot, and treasure icons from party card
  - [x] 23.1.5 unit test coverage in `internal/screens/hoard/hoard_test.go` and app integration in `cmd/pomquest/root.go`

## Feature 24: add party screen
- [x] 24.1 order dashboard quest, party, provisions, hoard, settings, quit
  - [x] 24.1.1 split provisions and party to party and provisions in `internal/screens/dashboard/dashboard.go`
  - [x] 24.1.2 wire `screenParty` navigation state in `cmd/pomquest/root.go`
- [x] 24.2 show full hero card stat blocks
  - [x] 24.2.1 support 10 RPG stat block view (`str`, `dex`, `cha`, `sur`, `con`, `wil`, `int`, `per`, `hel`, `luc`) when `ShowFullStats` is true in `internal/ui/herocard.go`
  - [x] 24.2.2 show active indicator icon in row below hero card
  - see example below:
```
╭─────────────╮ 
│pam         1│ 
│ranger       │ 
├─────────────┤ 
│str 0   dex 1│ 
│cha 0   sur 1│ 
│con 0   wil 0│ 
│int 0   per 1│ 
│hel 1   luc 0│ 
│             │ 
│󱡁 3🛡0    6/8│ 
╰─────────────╯ 
 󱐠              
```
- [x] 24.3 active hero card should have raised/hover effect by moving the card 2 col left and 1 col up
  - [x] 24.3.1 composite `compositeActiveOffsetCard2x1` to keep old border underneath while raising active card 2 cols left and 1 row up
  - [x] 24.3.2 create dedicated party screen in `internal/screens/party/party.go` with hero selection keybindings and unit tests in `party_test.go`
- [x] 24.4 render party based on sql data
  - [x] 24.4.1 show only number of hero cards based on the party max size
  - [x] 24.4.2 render empty hero slots as dotted border cards
  - [x] 24.4.3 inside of the empty card should be a bordered button that reads "Recruit Hero"
  - [x] 24.4.4 on selecting on enter'ing or r'ing should direct user to reroll screen
  


## Feature 25: add reroll screen
- [x] 25.1 a dead hero can be "rerolled" from the party screen
  - [x] 25.1.1 navigates to reroll screen when `r` or `enter` is pressed on a dead hero in `internal/screens/party/party.go`
- [x] 25.2 reroll screen looks very similar to party screen, showing a list of four hero cards
  - [x] 25.2.1 hero cards with no name
  - [x] 25.2.2 the first card is always "random" and has `?` for stats (`IsUnknownStats = true` in `internal/ui/herocard.go`)
  - [x] 25.2.3 the other three cards are 3 randomly selected classes from the class catalog programmatically
- [x] 25.3 reroll button
  - [x] 25.3.1 `[r]` rerolls available class options (cards 1, 2, 3), keeping card 0 ("random") untouched
  - [x] 25.3.2 only 2 rerolls allowed (`RerollsLeft` counter tracked and rendered in header)
- [x] 25.4 after hero picks a hero class card they pick a name
  - [x] 25.4.1 prompts for hero name input after card selection, replaces dead hero in party, and returns to party screen


## Feature 26: add provisions screen
Soand Sow's General Store
- [x] 26.1 create item catalog
  - [x] 26.1.1 contains items with base cost
    - [x] 24.1.1.1 food: 10c / lb
    - [x] 24.1.1.2 spare parts: 60c / part
    - [x] 24.1.1.3 ox: 20s per yoke (3 ox)
    - [x] 24.1.1.4 potions: 20s / pot
    - [x] 24.1.1.5 bombs: 50c / bomb
    - [x] 24.1.1.6 shovels: 20c / shovel
    - [x] 24.1.1.7 rope: 10c / 50ft
- [x] 26.2 add store screen
  - [x] 26.2.1 store props: name: Soand Sow's, type: general, sells: all items at base cost
- [x] 26.3 update coin value
  - [x] 26.3.1 conversion rate 1g = 10s; 1s = 100c
  - [x] 26.3.2 reduce database values to coin representation
  - [x] 26.3.3 gold and silver are visual representations of total coin value
- [x] 26.4 instead of owned show party card under store
  - [x] 26.4.1 party card is a shared ui component between provisions and quests
- [x] 26.5 set current card icons as general icon for each supply icon
  - [x] 26.5.1 defined once in icon package registry
- [x] 26.6 store items contain 3 cols (Item, Cost, Subtotal)
  - [x] 26.6.1 all columns aligned left
- [x] 26.7 when selecting item
  - [x] 26.7.1 prompt user for amount to purchase
  - [x] 26.7.2 update subtotal column dynamically
- [x] 26.8 last option is purchase
  - [x] 26.8.1 purchase row: purchase, blank cost, total subtotal
  - [x] 26.8.2 do not change party coin values until after purchase
- [x] 26.9 remove border from store










## Feature 27: create ascii image library
- [x] 27.1 use plan/ascii-art.md:575 for "treasure_hoard"
  - [x] 27.1.1 embed `TreasureHoardASCII` art in `internal/screens/hoard/hoard.go`

## Feature 28: connect heros and party to sqlite db
- [ ] 28.1 connect screens content for hero and party to sqlite db

## Feature 29: add settings screen
- [x] 29.1 contains
  - [x] 29.1.1 reset game data
    - [x] 27.1.1.1 resets and migrates the local sql db
  - [x] 29.1.2 audio
    - [x] 27.1.2.1 set in game audio - stub for now
  - [x] 29.1.3 credits - stub for nowt

## Feature 30: db flow
- [x] 30.1 implement db flow
  - [x] 30.1.1 first open -> init and migrate db 
  - [x] 30.1.2 second + open -> load db
  - [x] 30.1.3 reset db from settings as needed
- [x] 30.2 show party current and max size on dash board under options 
- [x] 30.3 show party coins to right of party hero count
- [x] 30.4 do not allow players with 0 current party members to "Quest!"

## Feature 31: create dungeon name generator
- [x] 31.1 create internal/dungeon/dungeon_name_generator
  - [x] 31.1.1 create book_of_location_names.go
  - [x] 31.1.2 generate random location name like donjon https://donjon.bin.sh/fantasy/name/#type=mon
  - [x] 31.1.3 general name generator util in internal/utils/name_generator.go
  - [x] 31.1.4 select template (locay of blank, blank's locay, blank locay)
  - [x] 31.1.5 select locay (dungeon | cave | cavern)
- [x] 31.2 generate name of dest: in quest screen based on template + generated name


