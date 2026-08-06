# Goal
 
pomodoro oregon trail like

oregon trail supplies

oxen
food
clothing
spare parts

# References

https://oregontrail.ws/games/the-oregon-trail/play/

# Architecture

Use existing template as base architecture.

Full page bubbletea tui view. 

GORM for sqllite db management.

Config.toml already in place for settings.

Use best practice memory management techniques to prevent memory leaks for long running usage. 
Clear in memory variables and timers when done with them. Use sqllite for data persistence and config.toml for settings management.

# Features

Gamified pomodoro time tracking with oregon trail and dnd flair. 
Persistent state across usage to progress in game goals.

### possible
- run in background ... or more rather seemless pickup. 
- notifications if running in background without being interacted with?


# Pitch

You are an ancient (dragon|dark lord) who is (trapped in a sacred mirror|ring) contrats with heros to do quests to increase your horde. 
  
Heroes are fragile replaceable things, but they have their quirks and desires. 

However you can not control the heroes, merely guide them on their journey.

# Pillars

slow deliberate progression                           
                                                      
quest > path/rest > encounter > interaction

all timers and counters use some fixed value + randomness + bonuses to calculate total time / success rate. 
                                                      
user dispatches or controls small party of 4-5 adventurers

# Sprites

Use nerd fonts for all icons. use ascii art for visuals

### dice
use d6 based system to use nerd font icons 󰇊 󰇋 󰇌 󰇍 󰇎 󰇏 
... could use ascii art for larger die rolls?  maybe v2? maybe d20's under the hood? or 2d6?
 
light d10
󰲠 󰲢 ... 󰿬

dark d10
󰲡 ... 󰿭
130 * 105 / 60 = 227.5
130 * 105 = 13650
### Creature sprites
ghost  
pacman, idk 󰮯 
blob 
sphagetti monster 
dog  
bug    
spider  󱇪 󱇫 󰯉
 󱖿  
󰭟

🧟‍♂


weapons
󱡂 󰓥 󰞇  󱡁 



### environment
blade trap 󱦗
(poison) spray trap 󱦝
󰠚   
 


### buildings
inn 
search house󱎰
house    󰹖     󰾁 󰀭 
healing 
󱡵 
 
### status sprites
fire 󰈸 
ice 
poison 
sleep 󰒲
dysentery 


# Heros
options for party cards

## Card ui
### detail view
            
 ╭───────────────╮
 │jimmothy       │─╮  ╭────────────╮  ╭─────────────╮  ╭─────────────╮    
 │fighter      1│ │  │applebee  2│  │pam        1│  │sam   (ded) │    
 ├───────────────┤ │  │wizard      │  │ranger       │  │rouge        │    
 │ str 2   dex 0 │ │  ├────────────┤  ├─────────────┤  ├─────────────┤    
 │ cha 0   sur 0 │ │  │str 0  dex 0│  │str 0   dex 1│  │str 0   dex 2│    
 │ con 2   wil 0 │ │  │cha 0  sur 0│  │cha 0   sur 1│  │cha 0   sur 0│    
 │ int 0   per 0 │ │  │con 0  wil 2│  │con 0   wil 0│  │con 0   wil 0│    
 │ hel 0   luc 0 │ │  │int 2  per 0│  │int 0   per 1│  │int 0   per 1│    
 │               │ │  │hel 2  luc 0│  │hel 1   luc 0│  │hel 0   luc 1│    
 │ ⚔ 2 🛡2   10 │ │  │            │  │             │  │             │    
 ╰───────────────╯ │  │ 3🛡0 6/6hp│  │󱡁 3🛡0    6/8│  │󰧼 3🛡0    ~0/8~│    
  ╰────────────────╯  ╰────────────╯  ╰─────────────╯  ╰─────────────╯    
  󰈸 󰒲                  󱕼              󱐠                


### questing view
 ╭────────────╮
 │jimmothy   2│─╮  ╭────────────╮  ╭────────────╮  ╭────────────╮    
 │fighter     │ │  │applebee   1│  │pam        1│  │~sam~       ~2~│    
 │            │ │  │wizard      │  │ranger      │  │~rouge~       │    
 │󰈸 󰒲 󱐠       │ │  │            │  │            │  │            │    
 │⚔ 2🛡2 10/10│ │  │ 󱕼         │  │󱐠           │  │            │    
 ╰────────────╯ │  │ 3🛡0   6/6│  │󱡁 3🛡0   6/8│  │~󰧼 3🛡0   0/8~│    
  ╰─────────────╯  ╰────────────╯  ╰────────────╯  ╰────────────╯    

(60 - 6 - 4) / 4 = 12.5
12*4 = 
                                              

## classes:
defined by 2-3 traits they are proficient in with a value of 4 points allocated into each stat. which weigh into which stats they go up by. 
stats are random, weighted heavily towards proficient stats

guard - con 3, str 1
fighter - str 2, con 1, dex 1
wizard - int 3, will 1
warlock - cha 3, int 1
rouge - dex 2, per 1, luc 1
ranger - dex 2, sur 2
ranger - dex 2, sur 2
cleric - hel 3, will 1
shaman - hel 2, sur 2
mage - int 4
bard - cha 4


### subclasses / personality
- 1 stat in a semi random attribute (fighter-healer) with 1 point in healing + 1 weight toward electing to level healing on level up. 

### perk
each class gets  special perk

guard: defends another unit x times per combat
fighter: extra attack x times per combat
wizard: fireball
warlock: charm spell
rouge: crit bonus
ranger: harvesting bonus
cleric: 1 revive per quest
shaman: 
mage:
bard: vicious mockery: apply weakness to all enemies 
# Party

╭───────────────────────────────────────────────────────────────────────╮
│ Party: 󰓥  󰱯 󰧼                                                        │
├───────────────────────────────────────────────────────────────────────┤
│  1g 60s 00c                  󱕼 sick: 1                               │
│ 󱐟 food: 400lbs                󰗑 diet: balanced                  󱕼 x1  │
│  spare parts: 2              󰓅 pace: fast                       x0  │
│  ox: 4                        regents: 8                      󰜦 x0  │
╰───────────────────────────────────────────────────────────────────────╯

╭───────────────────────────────────────────────────────────────────────╮
│ Party: 󰓥  󰱯 󰧼                                                        │
├───────────────────────────────────────────────────────────────────────┤
│  1g 60s 00c                  󱕼 sick: 1                               │
│ 󱐟 food: 400lbs                 diet: sparse                    󱕼 x1  │
│  spare parts: 2              󰓅 pace: fast  ◄─                   x0  │
│  ox: 4                        regents: 8   │                  󰜦 x0  │
╰──────────────────────────────────────────────│────────────────────────╯
                                               │           
                              may remove pace due to pomodoro component
                                     

# Events

- Semi random events based on grid.
- 󰇂 Encounters happen at occurrence of 1m * randomDelta(0.5)
- 󰇂 fights or encounter sub events and rolls happen at fixed interval


A "troop" of "12" "goblins" "ambush" the party. 
  dmg targets are semi random based party member constitution / health / class(?). 


Jimmothy fell into a trap and took 2 dmg
...
Applebee spotted and disarmed a trap!
...
Pam found spare part s


# Statuses

Affect rolls for dmg recieved / blocked


## Individual

- Tired
- Sick
- Hungery
- Happy
- Healthy
### Sickness

  - Dysentery
  - Infection
  - Madness
  - Cursed

## Interpersonal

- In love with
- Gratitude
- Attracted to
- Mad at 
- Upset with
- Related to
    - maybe some heroes come as a pair?
    - related heroes can not be attacted to or in love with relations


# Turn order

Random weighted toward higher dex

# combat

Basically an auto battler
Insert table of enemy target priorities here:


|enemy|target prio|atk|def|hp|
|---|---|---|---|---|
|bat|int|1|1|1|
|goblin|tank|1|0|2|
|boar|tank|1|0|2|
|raven|tank|1|0|2|
|raven|tank|1|0|2|
|feral boar|tank|1|0|2|
|werewolf|tank|2|1|4|
|woodland|tank| | | |
|swamp|dex| | | |


# Upkeep

> Things for player to have to spend gold on.

Each hero has a daily fee.
Food, supplies, regents, medical treatment, hiring new heroes. 

Heroes can be robbed during rests,                                                                                                                    ──
Food can go bad. 
Wagon breaks down. 
Horses / ox die. 
Equipment breaks. 
Party members lose / and sometimes win money gambling in town
Party members sneak food. 
...

# Crafting

Party uses regents to craft potions / poisons / etc.
Party members may craft on their own using regents, player may also craft. 


# Open Questions

## what happens on tpk?
options
-> slow down combat / interactions if party member dies to draw out time, slow down possible tpk to hault on 1 member left, always dies in last second? 
-> offer to restart quest loop. (likely will be during last work cycle)
-> revive perk?
-> intervene as arch dragon: burn x % of horde to revive x party members

## How do dates work?

12 month system? 4 seasons? 

existing dates / seasons? made up but sound similar (use same abbreviations (Jan... Jantarthic))? 

```sh
"Today marks the 120th day of Wintithe, tomorrow will be 001 Springtithe, 1000 R.O."
```

...scale up from dungeons to wars & use year system?
    ... or scale interaction from dungeon level to dispatch level (party quested for 3 years came back with 2 magical items 200 gold ...)

"Today marks the 120th day of Wintithe, tomorrow will be 001 Springtithe, 100 a.p."

`a.p.` after pomodoro
`p.p.` post pomodoro
`pom`  pomodoro

"The year is 0 a.p." --> 0 grows to rand number between 1 and 6000 ?



# Quests


| lvl | quest | biom | paths | encounters | 
| --- | --- | --- | --- | --- | 
| 1 | magic item hunt|travel through forest|forgage for shelf parts|find x|
|  | dungeon delve|travel through forest|forgage for shelf parts|find x|
|  |  forest delve|travel through forest|forgage for shelf parts|find x|
|  | | | | ... |
| 10 | destroy the ring|travel through forest|forgage for shelf parts|find x|
| | | | |ambushed by x|
| | | | |come across x|
| | | | |come across x|
| | save the kingdom|troll hunt|forgage for supplies|find x|
| | | |investigate sighting areas| | 
| | rescue the x|troll hunt|forgage for supplies|find x|
| | | |investigate sighting areas| | 




## Questing view 


 ╭─ 12, Jan, 1000 A.P. ───────────────────────────────────────────────╮
 │ ...                                                                │
 │ 6 goblins ambush the party!                                        │
 │ ...                                                                │
 │ jimmothy attacks 󰇍 󰇎 2/2 hit for ... 󰇊 󰇏 dmg!                      │
 │ ...                                                                │
 │ 4 goblins remaining...                                             │
 │ ...                                                                │
 │ applebee casts magic missle 󰇌 󰇎 󰇋  1/3 hit for ... 󰇏 /2 = 3 dmg!   │
 │ ...                                                                │
 │ last remaining goblin fled!                                        │
 │ ...                                                                │
 │ `+ 3 xp`! +3 copper!                                                 │
 ╰────────────────────────────────────────────────────────────────────╯
                                                       
     󰓥 7!
 ╭────────────╮                                                  
 │jimmothy    │─╮    ╭────────────╮    ╭────────────╮    ╭────────────╮    
 │fighter    2│ │    │applebee    │    │pam         │    │~sam~        │    
 │            │ │    │wizard     1│    │ranger     1│    │~rouge~      ~2~│    
 │󰈸 󰒲 󱐠       │ │    │            │    │            │    │            │    
 │⚔ 2🛡2 10/10│ │    │ 󱕼         │    │󱐠           │    │            │    
 ╰────────────╯ │    │ 3🛡0   6/6│    │󱡁 3🛡0   6/8│    │~󰧼 3🛡0   0/8~│    
  ╰─────────────╯    ╰────────────╯    ╰────────────╯    ╰────────────╯    
 ╭─────────────────────────────────────────────────────────────────────╮
 │ Party: 󰓥  󰱯 󰧼                                                 x0  │
 ├─────────────────────────────────────────────────────────────────────┤
 │  1g 60s 00c                  󱕼 sick: 1                             │
 │ 󱐟 food: 400lbs                󰗑 diet: balanced                󱡄 x1  │
 │  spare parts: 2              󰓅 pace: fast                     x0  │
 │  ox: 4                        regents: 8                    󰜦 x0  │
 ╰─────────────────────────────────────────────────────────────────────╯
  Dest: Forest of Izelliador                                21:00/25:00  
  ████████████████████████████████████████████████                    ] 
  Weather:   | Terrain:  | Hazards: mud, rain         days passed: 12




  --- event view ------------------------------------------------------
 ╭─ 22, Mar, 1023 A.O. ────────────────── Unamed dungeon near Perth ──╮ <- possible option for handling showing floors # and room # in dungeon / location? 
 ╭──────────────────────────────────────────────── Floor 1 | Room 2 ──╮     reverts back to fuller view internal areas aren't needed?
 │ ...                                                                │
 │ ...                                                                │
 │ jimmothy attacks 󰇍 󰇎 2/2 hit for ... 󰇊 󰇏 dmg!                      │
 │ ...                                                                │
 │ 4 goblins remaining...                                             │
 │ ...                                                                │
 │ applebee casts magic missle 󰇌 󰇎 󰇋  1/3 hit for ... (󰇏/3)+1= 3 dmg! │
 │ ...                                                                │
 │ last remaining goblin fled!                                        │
 │ ...                                                                │
 │ `+ 3 xp`! +3 copper!                                                 │
 ╰────────────────────────────────────────────────────────────────────╯
  --- dmg +/- view ----------------------------------------------------
     󰓥 7!
  --- party view ------------------------------------------------------
 ╭────────────╮                                                  
 │jimmothy    │─╮    ╭────────────╮    ╭────────────╮    ╭────────────╮    
 │fighter    2│ │    │applebee    │    │pam         │    │~sam~        │    
 │            │ │    │wizard     1│    │ranger     1│    │~rouge~      ~2~│    
 │󰈸       󰒲 󱐠 │ │    │            │    │            │    │            │    
 │⚔ 2🛡2 10/10│ │    │         󱕼 │    │          󱐠 │    │            │    
 ╰────────────╯ │    │ 3🛡0   6/6│    │󱡁 3🛡0   6/8│    │~󰧼 3🛡0   0/8~│    
  ╰─────────────╯    ╰────────────╯    ╰────────────╯    ╰────────────╯    
  |current turn |    |            |                      | dead       |
  --- party info ------------------------------------------------------
 ╭─────────────────────────────────────────────────────────────────────╮
 │ Party: 󰓥  󰱯 󰧼  <- turn order                        potions ->  2 │
 ├─────────────────────────────────────────────────────────────────────┤
 │  1g 60s 00c                  󱕼 sick: 1                             │
 │ 󱐟 food: 400lbs                󰗑 diet: balanced                 x0  │
 │  spare parts: 2              󰓅 pace: fast                    󱡄 x0  │
 │  ox: 4                        regents: 8                    󰜦 x1  │
 ╰─────────────────────────────────────────────────────────────────────╯
  --- quest info ------------------------------------------------------
  Dest: Forest of Izelliador                                21:00/25:00  
  ████████████████████████████████████████████████                    ] 
  Weather:   | Terrain:  | Hazards: mud, rain         days passed: 12


 ╭─ 12, Jan, 1000 A.P. ────────────────╮
 │ ...                                 │
 │ 6 goblins ambush the party!         │
 │ ...                                 │
 │ jimmothy attacks 󰇍 󰇎 2/2 hit for... │
 │ ... 󰇊 󰇏 dmg!                        │
 │ ...                                 │
 │ 4 goblins remaining...              │
 │ ...                                 │
 │ applebee casts magic missle ...     │
 │ ... 󰇌 󰇎 󰇋  1/3 hit for...           │
 │ ... 󰇏 /2 = 3 dmg!                   │
 │ ...                                 │
 │ last remaining goblin fled!         │
 │ ...                                 │
 │ `+ 3 xp`! +3 copper!                  │
 ╰─────────────────────────────────────╯
                                                   
                            󰓥 7!
  | 󰓥 Jimmothy          ╭────────────╮                                                  
                        │jimmothy    │─╮ 
   Pam                 │fighter    2│ │   
                        │            │ │     
  󰱯 Sam                 │󰈸 󰒲 󱐠       │ │ 
                        │⚔ 2🛡2 10/10│ │ 
  󰧼 Spam                ╰────────────╯ │   
                         ╰─────────────╯ 
 ╭─────────────────────────────────────╮
 │ Party: 󰓥  󰱯 󰧼                 x0  │
 ├─────────────────────────────────────┤
 │  1g 60s 00c   󱕼 1                  │
 │ 󱐟 400lbs       󰗑 balanced        󱡄 1│
 │  2            󰓅 fast             0│
 │  4             8               󰜦 0│
 ╰────────────────────────────────────────────────────────────╯
  Dest: Forest of Izelliador                        21:00/25:00  
  █████████████████████████████████████████████               ] 
    |  |  mud, rain   12

# Oregon

you may:
  0. travel the trail
  0. Learn about the trail
    - exposition
  3. See the top ten
  4. turn off sound
  5. choose mgmt options
    - see the current top 10
    - see the original top 10
    - erase the current top 10
    - erase umbsone messages
    - erase saved game
    - turn joystick on
    - calibrate
    - return 
  6. end

## travel the trail
  chose bank, carpenter, farmer. $$$ -> $ => less pts -> more pts
  shop general store
    oxen, food, cloting, ammo, spare parts - wheel, axle, tongue
    - "good luck"

## independence may 1, 1848
weather: cool
health: good
pace: steady
rations: filling
you may
  0. continue
  0. check supplies
  0. look at map
  0. change pace
  0. change food rations
  0. attempt to trade
  0. talk to people
  0. buy supplies

## continue
From __ it is __ miles to the __ River crossing.
date: may 1, 1884
weather: cool
health: good
food: 2000 pounds
next landmark: 102 miles
miles traveled: 0 miles

- press enter to size up the situation

## you are no at the ___ river crossing. Would you like to look around?
y -> "cutscene image"

## River crossing
> meta

0. continue
0. check sup 
0. look map
0. change pace
0. change rations
0. stop to rest
0. attempt to trade
8. talk to people

### crossing

> must cross the river to continue. 
> 637 ft across and 6 ft deep in the middle


weather: cool
River width: 637 feet
River depth: 6.0 feet
you may:

0. attempt to ford the river
0. caulk the wagon and float it across
3. take a ferry across
4. wait to see if coniditions improve
5. get more information

2.
-> wagon tipt lost 200lbs of food, 75 bullets, and ___ drowned (died)

### continue

"from the ___ River crossing it is __ miles to the Big Blue River crossing."
date: may 8
weather: hot
food: 1895 lbs
next landmark: 82 miles
miles traveled: 103 miles


### size up
 
dashboardy...

date
weather
health
pace
rations

0. continue on trail
0. check sup
0. look at map
0. change pace
0. change rations
0. stop to rest
0. attempt to trade
0. hunt for food

### random event
find wild fruit.

broken wagon tongue 
- try to repair it? 
  - fail. must use spare part. 
    - since you don't have a spare tongue, you must trade for one.
      you must trade for a tongue to be able to continue 
        -> spend days waiting for tongue to trade

### you are now at

you are now at fort ___. would you like to look around?

...
t
From fort -__ it is 250 miles to chimney rock.
dat: may 25, 1848
weather: cool
health: good
food: 1187 lbs
next land mark: 250 miles
miles traveled 304 miles

### random event
ox wanders off. lose 1 day.

you pass a gravesite. woud you like to look closer?

here lies a ...

# path options

dungeon -> dungeon -> dungeon

to dungeon -> dungeon -> dungeon -> back 
to dungeon -> dungeon -> back
forest -> dungeon -> deeper
forest -> dungeon -> castle
forest -> mountain -> river
swamp -> forest -> dungeon


