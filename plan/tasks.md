# Context

## Goal
The goal of this projects to create tui project oult

Concept layout more thoroughly in plan/mission.md

# Features
## Feature 1: define data models
- [ ] 1.1 define / refine gorm data models for game in plan/data_models.go.
- do not integrate into app yet. 
- data models will be used to represent everything from game state, to monsters, to heroes. 

  

## Feature 2: define sqlite db architecture
document:
- put findings in plan/sqlite-discovery.md
- what should be stored? What should be embedded enums? 
- define simple migration patterns

## Feature 3: define ui
document:
- put findings in plan/define-ui-discovery.md
- [ ] 3.1 create map of user flows
- based on plan/screens.md
- [ ] 3.2 tree diagram of ui architecture 
- should show proposed components
- each component should have a generic lib component and 1 or several domain specific components that utilize the lib components
  - to  optimize clarity and shared theming. 

## Feature 4: simple dashboard screen


## Feature 5: simple questing screen



## Feature 6: define class base stats
only warrior
str 2, con 2⌋
dmg: 1d6/2+1 => 0.5, 2, 2.5, 3, 3.5, 4
dmg: 1d6 % 4 = 2dmg! 
dmg: (1d6+2)/2 = 1.5, 2, 2.5, 3, 3.5, 4
dmg: ⌊(1d6+2)/2⌋ = 1, 2, 2, 3, 3, 4

maybe this
dmg: hits for 1d2 󰇋 (2) dmg!
## Feature 7: create randomization library
2d6 based system 2.6% crit chance - epic fail


## Feature 8: define quest, biome, path, weather tables
start with 1 quest, 1 biome, 1 path, 1 sunny weather


## Feature 9: define monster stat table

... maybe sql isn't needed for stats ext? unless, they are downloaded with first install... sqlite still good for saving player state. But does it make sense to ship all as embeded sqlite db?

... all stats could really be enums. what needs to saved is history, and player state... that could be file... sqlite would be good for large history storage, viewing previous runs, locations impacting each other, or revisited? 
- especially if using something like faker for names.


start with 1 monster, bats... 1hp 0def 1 atk
