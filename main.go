package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	_ "github.com/charmbracelet/bubbletea"
)

var version = "dev"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "txs",
		Short:   "Termux Station — single-binary TUI dashboard",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "txs version "+version)
			return nil
		},
	}
	return cmd
}
