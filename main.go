package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/Pastalikek65/txs/store"
	"github.com/Pastalikek65/txs/tui"
)

var version = "0.1.0"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var showVersion bool
	cmd := &cobra.Command{
		Use:     "txs",
		Short:   "Termux Station — single-binary TUI dashboard for Termux",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				fmt.Fprintln(cmd.OutOrStdout(), "txs version "+version)
				return nil
			}
			// plain fallback for CI
			if !term.IsTerminal(int(os.Stdout.Fd())) || os.Getenv("TXS_PLAIN")=="1" {
				fmt.Fprintln(cmd.OutOrStdout(), "txs 5 panels: github jobs db sys files — run in Termux for TUI")
				return nil
			}
			home, _ := os.UserHomeDir()
			if home=="" { home = os.Getenv("HOME") }
			if home=="" { home = "/tmp" }
			dbPath := filepath.Join(home, ".cache", "txs", "txs.db")
			if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
				dbPath = filepath.Join(xdg, "txs", "txs.db")
			}
			s, err := store.Open(dbPath)
			if err != nil { return err }
			defer s.Close()
			m := tui.NewRoot(s)
			p := tea.NewProgram(m, tea.WithAltScreen())
			_, err = p.Run()
			return err
		},
	}
	cmd.Flags().BoolVarP(&showVersion, "version", "v", false, "show version")
	return cmd
}
