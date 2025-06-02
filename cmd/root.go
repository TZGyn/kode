/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/TZGyn/kode/internal/model"

	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	huh "github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "kode",
	Short:         "CLI AI Assistant",
	Long:          "CLI AI Assistant",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		checkGitCmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
		stdout, err := checkGitCmd.Output()

		if err != nil {
			return errors.New("invalid git repo")
		}

		if strings.Split(string(stdout), "\n")[0] != "true" {
			return errors.New("invalid git repo")
		}

		var c model.ChatConfig
		configFilePath, err := xdg.ConfigFile(filepath.Join("kode", "kode.json"))
		if err != nil {
			return err
		}

		dir := filepath.Dir(configFilePath)
		if err = os.MkdirAll(dir, 0o700); err != nil {
			return err
		}

		if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {

			f, err := os.Create(configFilePath)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			defaultConfig, _ := json.MarshalIndent(&c, "", "\t")
			_, err = f.WriteString(string(defaultConfig))
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		content, err := os.ReadFile(configFilePath)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(content, &c); err != nil {
			return err
		}
		opts := []tea.ProgramOption{}
		opts = append(opts, tea.WithOutput(os.Stderr))

		for {
			var prompt string

			promptForm :=
				huh.NewText().Title("Enter a prompt:").
					Value(&prompt)

			err = promptForm.Run()

			if err != nil && err == huh.ErrUserAborted {
				return errors.New("user canceled")
			} else if err != nil {
				return errors.New("prompt failed")
			}

			out, err := glamour.Render(prompt, "auto")
			if err != nil {
				fmt.Println(err)
				break
			}

			fmt.Print("\n  ")
			color.BgRGB(95, 95, 255).AddRGB(230, 255, 219).Println(" User: ")
			fmt.Println(out)

			chatModel := model.InitialModel(prompt, c)

			p := tea.NewProgram(chatModel, opts...)
			m, err := p.Run()
			if err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				break
			}

			chatModel = m.(*model.ChatModel)

			if chatModel.Response != "" {
				out, err := glamour.Render(chatModel.Response, "auto")
				if err != nil {
					fmt.Println(err)
				}
				fmt.Print("  ")
				color.RGB(254, 229, 250).AddBgRGB(245, 127, 224).Println(" Assistant: ")
				fmt.Println(out)
			} else {
				fmt.Println("No Response")
			}
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kode.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
