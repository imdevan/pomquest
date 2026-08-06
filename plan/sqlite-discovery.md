# SQLite Database Architecture & Strategy Discovery

## 1. Data Storage Breakdown: SQLite DB vs Embedded Enums & Templates

### What Should Be Stored in SQLite DB (Dynamic Persistent State)
The SQLite database stores all **runtime state**, **user progression**, and **historical records** that persist across app launches and Pomodoro sessions:

- **Active Session State (`GameState`)**: Active quest ID, active party ID, active location ID, active pomodoro timer status, total focus time, current loop count, and active session flag.
- **Party & Hero Dynamics (`Party`, `Hero`)**: Current HP, max HP, status ailments (e.g. Dysentery, Tired, Sick), dynamic stats, experience, level, items, gold, food (lbs), spare parts, oxen, regents, potions, poisons.
- **Quest & Encounter Records (`Quest`, `Encounter`)**: Current quest progress, distance covered, weather state, encounter log outcomes.
- **Global Meta Progression (`Player`, `Hoard`)**: Dragon/Dark Lord hoard metrics (chests, magic items, total gold collected across runs), cumulative runs, heroes recruited vs lost.
- **Game History Logs (`GameLog`)**: Chronological event stream messages for UI scrollback and session recap.

### What Should Be Embedded Enums & In-Code Templates
Static rules, definitions, and archetypes should **NOT** clutter SQLite tables as initial data seeds. Instead, they live in Go code as typed constants/struct templates (or embedded via `go:embed`):

- **Enums & State Flags**:
  - `PomodoroState`: `Idle`, `Working`, `ShortBreak`, `LongBreak`
  - `HeroStatus`: `Healthy`, `Tired`, `Sick`, `Dysentery`, `Cursed`, `Dead`
  - `DietPace`: `Sparse`, `Balanced`, `Generous`
  - `TargetPriority`: `Tank`, `Int`, `Dex`, `LowestHP`
  - `WeatherType`: `Sunny`, `Rain`, `Storm`, `Fog`
- **Archetype Templates**:
  - `ClassDef` templates (Fighter: Str 2/Con 1/Dex 1; Guard: Con 3/Str 1; Wizard: Int 3/Wil 1).
  - Monster base templates (Bat: 1 HP, 0 Def, 1 Atk; Goblin: 2 HP, 0 Def, 1 Atk).
  - Biome & Weather modifiers.

*Why?* Keeping templates in Go code keeps migrations trivial, avoids database seed corruption, provides compile-time safety, and simplifies unit tests.

---

## 2. Migration Pattern & Database Location

### File Location
Following the XDG Base Directory Specification:
- **Linux/macOS**: `$XDG_DATA_HOME/pomquest/pomquest.db` (defaults to `~/.local/share/pomquest/pomquest.db`).
- Configured via `internal/config`.

### Migration Pattern
- **GORM AutoMigrate**: Call `db.AutoMigrate(...)` on application initialization for all models (`Hero`, `ClassDef`, `Monster`, `NPC`, `Party`, `Hoard`, `Quest`, `Location`, `Encounter`, `GameState`, `Player`, `GameLog`).
- **Pragmas**:
  - `PRAGMA journal_mode = WAL;` (Write-Ahead Logging for fast concurrent reads & durability).
  - `PRAGMA foreign_keys = ON;` (Enforce referential integrity).
  - `PRAGMA synchronous = NORMAL;` (Optimal balance of safety and speed for CLI).
- **Schema Version Tracking**: Maintain a lightweight `schema_migrations` table or GORM model version check for future destructive schema changes requiring custom transformation functions.

---

## 3. Application Interaction Pattern

### Architecture Flow
```
+-------------------------------------------------------+
|                 Bubble Tea TUI Views                  |
+-------------------------------------------------------+
                           |
                           v
+-------------------------------------------------------+
|            Domain Service / Game Engine               |
|      (Pomodoro Logic, Combat Engine, Upkeep)         |
+-------------------------------------------------------+
                           |
                           v
+-------------------------------------------------------+
|             GameRepository Interface                  |
+-------------------------------------------------------+
                           |
                           v
+-------------------------------------------------------+
|            GORM / SQLite Persistence Layer            |
|              (Max 1 Open DB Connection)              |
+-------------------------------------------------------+
```

### Connection Management
- **Single Writer Connection Pool**: Set `db.SetMaxOpenConns(1)` to eliminate `SQLITE_BUSY` database lock conflicts in CLI environments.
- **Repository Interface**: TUI components interact with `GameRepository` (e.g., `SaveGameState()`, `LoadActiveSession()`, `LogEvent()`) rather than executing raw SQL queries directly.

---

## 4. Handling Early Session Termination (Crashes, Ctrl+C, Terminal Close)

### Auto-Checkpointing Strategy
- **Frequent Micro-Saves**: Persist `GameState`, `Party` supplies, and `Hero` health after every key event:
  - Pomodoro timer minute ticks.
  - Combat / encounter completion.
  - Party upkeep calculation (food consumption, rests).
- **WAL Mode Transaction Safety**: SQLite Write-Ahead Logging guarantees ACID compliance. An unexpected terminal crash (`SIGKILL` or power failure) will never corrupt the database file; uncommitted transactions revert cleanly.

### Crash Recovery & Session Resume
- On app launch, query for `GameState` with `IsActive = true`.
- If an active session is found:
  1. Calculate wall-clock elapsed time: `elapsed := time.Since(gameState.UpdatedAt)`.
  2. If the user closed during a working pomodoro:
     - Prompt user in TUI: **"Resume active session"** or **"Abandon quest & rest"**.
     - If resumed, adjust remaining timer by elapsed time or pause at last clean checkpoint.

---

## 5. Key Additional Architectural Questions & Considerations

### Q1: How to handle multiple terminal instances running `pomquest` simultaneously?
- **Solution**: Use file locking (`pomquest.lock`) in the app data directory or configure a GORM busy timeout (`PRAGMA busy_timeout = 5000;`). If another instance is running, display a helpful message: *"PomQuest is already running in another terminal."*

### Q2: How do we prevent `GameLog` from bloating memory and disk over long usage?
- **Solution**: Implement log pruning / rotation in `GameRepository`. Retain the last 1,000 log entries per session and automatically archive or prune older entries.

### Q3: How will unit testing work without touching disk?
- **Solution**: Use SQLite in-memory connections (`sqlite.Open("file::memory:?cache=shared")`) in unit and integration tests.

### Q4: How can users manage or backup their save data?
- **Solution**: Provide CLI commands:
  - `pomquest db path` - Print active database location.
  - `pomquest db reset` - Reset game state / re-initialize database.
