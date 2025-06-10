/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		configFilePath, err := xdg.ConfigFile(filepath.Join("kode", "kode.json"))
		if err != nil {
			return err
		}
		operatingSystem := runtime.GOOS

		var openCmd *exec.Cmd
		switch operatingSystem {
		case "windows":
			openCmd = exec.Command(filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe"), "url.dll,FileProtocolHandler", configFilePath)
		case "darwin":
			openCmd = exec.Command("open", configFilePath)
		case "linux":
			openCmd = exec.Command("xdg-open", configFilePath)
		}

		if openCmd == nil {
			return errors.New("cannot open file")
		}

		err = openCmd.Run()
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
