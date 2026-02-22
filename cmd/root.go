/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/TZGyn/kode/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:           "kode",
	Short:         "CLI AI Assistant",
	Long:          "CLI AI Assistant",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(
			model.InitAppModel(),
			tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
			tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
		)
		if _, err := p.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}
