// contagious list <reg> -n -o matrix.json
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/duffney/contagious/internal/ghcr"
	"github.com/duffney/contagious/tagStore"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Args:  cobra.ExactArgs(1),
	Long:  ``,
	Run:   handleList,
}

var outputFile string

func init() {
	listCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
	listCmd.Flags().BoolP("next-tag", "n", false, "Include next patch tag information")
}

func handleList(cmd *cobra.Command, args []string) {
	var username string
	registry := args[0]
	tagStore := tagStore.New()
	nTag, _ := cmd.Flags().GetBool("next-tag")

	if strings.Contains(registry, "ghcr.io") {
		username = strings.Split(registry, "/")[1]
	}

	images, err := ghcr.ListImages(username)
	if err != nil {
		os.Exit(1)
	}

	for _, img := range images {
		tagStore.AddImage(img)
	}

	if outputFile == "" {
		out, _ := tagStore.PrintTable(nTag)
		fmt.Println(out)
		return
	}

	jsonOutput, _ := tagStore.GetJSON(nTag)
	err = os.WriteFile(outputFile, []byte(jsonOutput), 0644)
	if err != nil {
		os.Exit(1)
	}
}
