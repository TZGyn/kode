/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/TZGyn/kode/internal/config"
	"github.com/TZGyn/kode/internal/message"
	"github.com/TZGyn/kode/internal/model"
	"github.com/TZGyn/kode/internal/models"

	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	huh "github.com/charmbracelet/huh"
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

		c, err := config.New()
		if err != nil {
			return nil
		}

		opts := []tea.ProgramOption{}
		opts = append(opts, tea.WithOutput(os.Stderr))

		messages := model.ChatMessages{}

		providerOpts := make([]huh.Option[models.ModelProvider], 0, len(models.Models))
		modelOpts := map[models.ModelProvider][]huh.Option[models.ModelID]{}

		for provider, models := range models.Models {
			providerOpts = append(providerOpts, huh.NewOption(string(provider), provider))
			for _, model := range models {
				modelOpts[provider] = append(modelOpts[provider], huh.NewOption(string(model.ID), model.ID))
			}
		}

		if c.DEFAULT_MODEL == "" || c.DEFAULT_PROVIDER == "" {
			err = huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[models.ModelProvider]().
						Title("Choose the Provider:").
						Options(providerOpts...).
						Value(&c.DEFAULT_PROVIDER),
					huh.NewSelect[models.ModelID]().
						TitleFunc(func() string {
							return fmt.Sprintf("Choose the model for '%s':", c.DEFAULT_PROVIDER)
						}, &c.DEFAULT_PROVIDER).
						OptionsFunc(func() []huh.Option[models.ModelID] {
							return modelOpts[c.DEFAULT_PROVIDER]
						}, &c.DEFAULT_PROVIDER).
						Value(&c.DEFAULT_MODEL),
				),
			).Run()

			if err != nil {
				return err
			}

			err = c.SaveConfig()
			if err != nil {
				return err
			}
		}

		for {
			var prompt string

			promptForm := huh.NewForm(
				huh.NewGroup(
					huh.NewText().Title("Enter a prompt:").
						Value(&prompt).Description("/model to update model"),
				),
			)

			err = promptForm.Run()

			if err != nil && err == huh.ErrUserAborted {
				return errors.New("user canceled")
			} else if err != nil {
				return errors.New("prompt failed")
			}

			if prompt == "/model" {
				err = huh.NewForm(
					huh.NewGroup(
						huh.NewSelect[models.ModelProvider]().
							Title("Choose the Provider:").
							Options(providerOpts...).
							Value(&c.DEFAULT_PROVIDER),
						huh.NewSelect[models.ModelID]().
							TitleFunc(func() string {
								return fmt.Sprintf("Choose the model for '%s':", c.DEFAULT_PROVIDER)
							}, &c.DEFAULT_PROVIDER).
							OptionsFunc(func() []huh.Option[models.ModelID] {
								return modelOpts[c.DEFAULT_PROVIDER]
							}, &c.DEFAULT_PROVIDER).
							Value(&c.DEFAULT_MODEL),
					),
				).Run()

				if err != nil {
					return err
				}

				err = c.SaveConfig()
				if err != nil {
					return err
				}
				continue
			}

			out, err := glamour.Render(prompt, "auto")
			if err != nil {
				fmt.Println(err)
				break
			}

			fmt.Println(message.UserStyle.Render(out))

			chatModel := model.InitialModel(prompt, messages, model.ChatConfig{
				Provider:          string(c.DEFAULT_PROVIDER),
				Model:             string(c.DEFAULT_MODEL),
				GEMINI_API_KEY:    c.GEMINI_API_KEY,
				OPENAI_API_KEY:    c.OPENAI_API_KEY,
				ANTHROPIC_API_KEY: c.ANTHROPIC_API_KEY,
			})

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
				fmt.Println(message.AssistantStyle.Render(out + "\n\n" + "  " + message.SecondaryStyle.Render(chatModel.Provider+" "+chatModel.Model)))
			} else {
				fmt.Println("No Response")
			}

			if len(chatModel.GoogleClient.Messages) > 0 {
				messages = model.ChatMessages{}
				err = messages.AddGoogleMessages(chatModel.GoogleClient.Messages)
				if err != nil {
					fmt.Println("Failed to remember google response")
				}
			}
			if len(chatModel.OpenAIClient.Messages) > 0 {
				messages = model.ChatMessages{}
				err = messages.AddOpenAIMessages(chatModel.OpenAIClient.Messages)
				if err != nil {
					fmt.Println("Failed to remember openai response")
				}
			}
			if len(chatModel.AnthropicClient.Messages) > 0 {
				messages = model.ChatMessages{}
				err = messages.AddAnthropicMessages(chatModel.AnthropicClient.Messages)
				if err != nil {
					fmt.Println("Failed to remember openai response")
				}

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
