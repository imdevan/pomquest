package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/pomquest/internal/config"
	"github.com/pomquest/internal/domain"
	pkg "github.com/pomquest/internal/package"
	"github.com/pomquest/internal/screens/dashboard"
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

	return runDashboard(cfg)
}

func runDashboard(cfg domain.Config) error {
	model := dashboard.NewModel(cfg)

	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithoutSignalHandler())
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run dashboard: %w", err)
	}

	if dbScreen, ok := finalModel.(dashboard.Model); ok {
		switch dbScreen.Selected {
		case dashboard.ChoiceQuest:
			fmt.Println("⚔ Starting a new Quest Journey...")
		case dashboard.ChoiceHoard:
			fmt.Println("💎 Dragon Hoard: 0 Gold, 0 Magic Items")
		case dashboard.ChoiceProvision:
			fmt.Println("🎒 Party Provisions: 100 lbs Food, 4 Oxen, 2 Spare Parts")
		case dashboard.ChoiceSettings:
			fmt.Println("⚙ Settings / Config options")
		case dashboard.ChoiceQuit:
			fmt.Println("Goodbye, Dark Lord!")
		}
	}

	return nil
}
