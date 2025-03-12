package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Displays the current version of the healthcheck CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("copamatic version %s\n", Version)
	},
}
