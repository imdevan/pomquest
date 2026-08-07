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

## Feature 6: party card
- [x] 6.1 add team hearts ui component
  - [x] 6.1.1 takes max and current hp, renders hearts evaluated in 10 HP and 5 HP increments using icon.KeyHeart...
  - [x] 6.1.2 add team hearts to party texture card top-right

## Feature 7: create custom internal/ui/texturedCard
- [x] 7.1 create custom internal/ui/texturedCard component supporting embedded text on any of the 4 border corners with padding, margin, and styling properties
- [x] 7.2 use TexturedCard for party card container wrapping party table



## Feature 8: party model update
- [x] 8.1 add these to the party model:
  - [x] 8.1.1 small non-negative ints (loot, rope, bombs, potions, shovels)
  - [x] 8.1.2 enum / string with fixed set of valid options (item usage: "none", "sparing", "balanced", "generous"; greed: "none", "low", "med", "high", "extreme")
- [x] 8.2 add those properties to the partycard

## Feature 9: create icon package
- [x] 9.1 map strings to nerdfont utf equivalents
- [x] 9.2 some icons can have a nerd font and emoji variant
- [x] 9.3 create library of names mapped to icons
- [x] 9.4 replace all icons on quest screen components with icon package icons

## Feature 10: create icon library

## Feature 11: define class base stats
only warrior
str 2, con 2⌋
dmg: 1d6/2+1 => 0.5, 2, 2.5, 3, 3.5, 4
dmg: 1d6 % 4 = 2dmg! 
dmg: (1d6+2)/2 = 1.5, 2, 2.5, 3, 3.5, 4
dmg: ⌊(1d6+2)/2⌋ = 1, 2, 2, 3, 3, 4

maybe this
dmg: hits for 1d2 󰇋 (2) dmg!
## Feature 12: create randomization library
2d6 based system 2.6% crit chance - epic fail


## Feature 13: define quest, biome, path, weather tables
start with 1 quest, 1 biome, 1 path, 1 sunny weather


## Feature 14: define monster stat table

... maybe sql isn't needed for stats ext? unless, they are downloaded with first install... sqlite still good for saving player state. But does it make sense to ship all as embeded sqlite db?

... all stats could really be enums. what needs to saved is history, and player state... that could be file... sqlite would be good for large history storage, viewing previous runs, locations impacting each other, or revisited? 
- especially if using something like faker for names.


start with 1 monster, bats... 1hp 0def 1 atk

## Feature 15: add icons
- [x] 15.1 add heart icons in internal/adapters/icon/icon.go:74    
  - [x] 15.1.1 10hp heart 󰋑 
  - [x] 15.1.2 10hp - 5 heart 󰛞 
  - [x] 15.1.3 10hp - 10 heart 󰋕
  - [x] 15.1.4 5hp heart 󰛟
  - [x] 15.1.5 5hp empty heart 󰛠


