package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "contagious",
	Short: "A CLI that lists patchable tags for Copacetic",
	Long:  `Contagious is a command line tool that generates a .json file listing container images tags for Copacetic.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listCmd)
}
