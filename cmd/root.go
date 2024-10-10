package cmd

import (
	"os"

	"conway-life-go/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "conway-life-go",
	Short: "Simple implementation of Conway's Game of Life in Golang with a TUI",
	RunE: func(_cmd *cobra.Command, _args []string) error {
		return tui.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.conway-life-go.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
