package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/pomquest/internal/app"
	"github.com/pomquest/internal/config"
	"github.com/pomquest/internal/domain"
	pkg "github.com/pomquest/internal/package"
	"github.com/pomquest/internal/screens/dashboard"
	"github.com/pomquest/internal/screens/hoard"
	"github.com/pomquest/internal/screens/party"
	"github.com/pomquest/internal/screens/questing"
	"github.com/pomquest/internal/screens/reroll"
	"github.com/pomquest/internal/screens/settings"
)

// Metadata loaded from package.toml at build time
var (
	version = pkg.Version()
	name    = pkg.Name()
	short   = pkg.Short()
)

type rootOptions struct {
	configPath  string
	showVersion bool
}

var rootCmd = newRootCmd()

// Execute is the CLI entrypoint.
func Execute() error {
	return rootCmd.Execute()
}

func newRootCmd() *cobra.Command {
	opts := &rootOptions{}
	cmd := &cobra.Command{
		Use:   name,
		Short: short,
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.showVersion {
				ver := resolvedVersion()
				cmd.Printf("%s\n", ver)
				return nil
			}
			return runInteractive(cmd, opts, args)
		},
	}

	cmd.Flags().StringVarP(&opts.configPath, "config", "c", "", "config file path")
	cmd.Flags().BoolVarP(&opts.showVersion, "version", "v", false, "print version information")

	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newCompletionCmd())

	return cmd
}

func resolvedVersion() string {
	ver := version
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ver
	}
	if ver == "dev" && strings.TrimSpace(info.Main.Version) != "" && info.Main.Version != "(devel)" {
		ver = info.Main.Version
	}
	return ver
}

func runInteractive(cmd *cobra.Command, opts *rootOptions, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	manager := config.NewManager(cwd)
	var cfg domain.Config
	if opts.configPath != "" {
		cfg, err = manager.LoadWithOverride(opts.configPath)
	} else {
		cfg, err = manager.Load()
	}
	if err != nil {
		cfg = domain.DefaultConfig()
	}

	model := newAppModel(cfg)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithoutSignalHandler())
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run application: %w", err)
	}

	if app, ok := finalModel.(appModel); ok && app.exitMessage != "" {
		fmt.Println(app.exitMessage)
	}

	return nil
}

type screenState int

const (
	screenDashboard screenState = iota
	screenQuesting
	screenParty
	screenReroll
	screenHoard
	screenSettings
)

type appModel struct {
	cfg         domain.Config
	dbPath      string
	state       screenState
	dashboard   dashboard.Model
	questing    questing.Model
	party       party.Model
	reroll      reroll.Model
	hoard       hoard.Model
	settings    settings.Model
	width       int
	height      int
	exitMessage string
}

func newAppModel(cfg domain.Config) appModel {
	return newAppModelWithDB(cfg, app.DefaultDBPath())
}

func newAppModelWithDB(cfg domain.Config, dbPath string) appModel {
	db, err := app.InitDB(dbPath)
	var activeParty app.Party
	if err == nil && db != nil {
		activeParty, _ = app.LoadOrCreateActiveParty(db)
	} else {
		activeParty = app.Party{
			Name:             "Adventuring Party",
			CurrentPartySize: 0,
			MaxPartySize:     1,
			Gold:             1,
		}
	}

	dbModel := dashboard.NewModelFromParty(cfg, activeParty)
	partyModel := party.NewModelFromParty(cfg, activeParty)

	return appModel{
		cfg:       cfg,
		dbPath:    dbPath,
		state:     screenDashboard,
		dashboard: dbModel,
		questing:  questing.NewModel(cfg),
		party:     partyModel,
		reroll:    reroll.NewModel(cfg, 0),
		hoard:     hoard.NewModel(cfg),
		settings:  settings.NewModelWithDBPath(cfg, dbPath),
		width:     80,
		height:    24,
	}
}

func (m appModel) returnToDashboard() (tea.Model, tea.Cmd) {
	m.state = screenDashboard
	var activeParty app.Party
	if m.dbPath != "" {
		if db, err := app.InitDB(m.dbPath); err == nil && db != nil {
			activeParty, _ = app.LoadOrCreateActiveParty(db)
		}
	}
	m.dashboard = dashboard.NewModelFromParty(m.cfg, activeParty)
	m.dashboard.SetDimensions(m.width, m.height)
	return m, m.dashboard.Init()
}

func (m appModel) Init() tea.Cmd {
	return m.dashboard.Init()
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.dashboard.SetDimensions(msg.Width, msg.Height)
		m.questing.SetDimensions(msg.Width, msg.Height)
		m.party.SetDimensions(msg.Width, msg.Height)
		m.reroll.SetDimensions(msg.Width, msg.Height)
		m.hoard.SetDimensions(msg.Width, msg.Height)
		m.settings.SetDimensions(msg.Width, msg.Height)
	}

	switch m.state {
	case screenDashboard:
		newDb, cmd := m.dashboard.Update(msg)
		if dbModel, ok := newDb.(dashboard.Model); ok {
			m.dashboard = dbModel
			if dbModel.Selected >= 0 {
				switch dbModel.Selected {
				case dashboard.ChoiceQuest:
					m.state = screenQuesting
					m.questing = questing.NewModel(m.cfg)
					m.questing.SetDimensions(m.width, m.height)
					return m, m.questing.Init()
				case dashboard.ChoiceParty:
					m.state = screenParty
					if db, err := app.InitDB(m.dbPath); err == nil && db != nil {
						if activeParty, err := app.LoadOrCreateActiveParty(db); err == nil {
							m.party = party.NewModelFromParty(m.cfg, activeParty)
						}
					}
					m.party.SetDimensions(m.width, m.height)
					return m, m.party.Init()
				case dashboard.ChoiceProvision:
					m.exitMessage = "🎒 Provisions: 100 lbs Food, 4 Oxen, 2 Spare Parts"
					return m, tea.Quit
				case dashboard.ChoiceHoard:
					m.state = screenHoard
					m.hoard = hoard.NewModel(m.cfg)
					m.hoard.SetDimensions(m.width, m.height)
					return m, m.hoard.Init()
				case dashboard.ChoiceSettings:
					m.state = screenSettings
					m.settings = settings.NewModelWithDBPath(m.cfg, m.dbPath)
					m.settings.SetDimensions(m.width, m.height)
					return m, m.settings.Init()
				case dashboard.ChoiceQuit:
					m.exitMessage = "Goodbye, Dark Lord!"
					return m, tea.Quit
				}
			}
		}
		return m, cmd

	case screenQuesting:
		newQ, cmd := m.questing.Update(msg)
		if qModel, ok := newQ.(questing.Model); ok {
			m.questing = qModel
			if qModel.Quitted {
				return m.returnToDashboard()
			}
		}
		return m, cmd

	case screenParty:
		newP, cmd := m.party.Update(msg)
		if pModel, ok := newP.(party.Model); ok {
			m.party = pModel
			if pModel.NavToReroll {
				m.party.NavToReroll = false
				m.state = screenReroll
				m.reroll = reroll.NewModel(m.cfg, pModel.RerollTargetIndex)
				m.reroll.SetDimensions(m.width, m.height)
				return m, m.reroll.Init()
			}
			if pModel.Quitted {
				return m.returnToDashboard()
			}
		}
		return m, cmd

	case screenReroll:
		newR, cmd := m.reroll.Update(msg)
		if rModel, ok := newR.(reroll.Model); ok {
			m.reroll = rModel
			if rModel.Completed {
				m.party.ReplaceHero(rModel.TargetHeroIdx, rModel.NewHero)
				m.state = screenParty
				m.party.SetDimensions(m.width, m.height)
				return m, m.party.Init()
			}
			if rModel.Quitted {
				m.state = screenParty
				m.party.SetDimensions(m.width, m.height)
				return m, m.party.Init()
			}
		}
		return m, cmd

	case screenHoard:
		newH, cmd := m.hoard.Update(msg)
		if hModel, ok := newH.(hoard.Model); ok {
			m.hoard = hModel
			if hModel.Quitted {
				return m.returnToDashboard()
			}
		}
		return m, cmd

	case screenSettings:
		newS, cmd := m.settings.Update(msg)
		if sModel, ok := newS.(settings.Model); ok {
			m.settings = sModel
			if sModel.Quitted {
				return m.returnToDashboard()
			}
		}
		return m, cmd
	}

	return m, nil
}

func (m appModel) View() string {
	switch m.state {
	case screenDashboard:
		return m.dashboard.View()
	case screenQuesting:
		return m.questing.View()
	case screenParty:
		return m.party.View()
	case screenReroll:
		return m.reroll.View()
	case screenHoard:
		return m.hoard.View()
	case screenSettings:
		return m.settings.View()
	default:
		return ""
	}
}
