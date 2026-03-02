/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/TZGyn/kode/internal/model"
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
			tea.WithContext(cmd.Context()),
		)
		if _, err := p.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}
