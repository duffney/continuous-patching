package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/duffney/contagious/internal/ghcr"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "contagious",
	Short: "A CLI that lists patchable tags for Copacetic",
	Long:  `Contagious is a command line tool that generates a .json file listing container images tags for Copacetic.`,
	Args:  cobra.ExactArgs(1),
	Run:   handleRoot,
}

var outputFile string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().BoolP("list", "l", false, "List patchable images from registry")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
}

func handleRoot(cmd *cobra.Command, args []string) {
	var userName string
	registry := args[0]

	// TODO: add support for dockerhub and acr
	if strings.Contains(registry, "ghcr.io") {
		userName = strings.Split(registry, "/")[1]
	}

	listFlag, _ := cmd.Flags().GetBool("list")
	if listFlag {
		patchTags := getPatchableTags(userName)

		if outputFile != "" {
			jsonOutput, err := json.MarshalIndent(patchTags, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			err = os.WriteFile(outputFile, jsonOutput, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to %s: %v\n", outputFile, err)
				os.Exit(1)
			}
		} else {
			_, err := os.Stdout.Write([]byte(strings.Join(patchTags, "\n")))
			if err != nil {
				os.Exit(1)
			}
		}
	}
}

func getPatchableTags(username string) []string {
	//TODO: replace with interface
	images, err := ghcr.ListImages(username)
	if err != nil {
		os.Exit(1)
	}

	data := make(map[string]int)

	for _, image := range images {
		name, tag, _ := strings.Cut(image, ":")

		if regexp.MustCompile(`-\d+$`).MatchString(tag) {
			patchIndex := strings.LastIndex(tag, "-")
			patchCounter, _ := strconv.Atoi(tag[patchIndex+1:])
			imageWithoutPatch := name + ":" + tag[:patchIndex]

			if value, ok := data[imageWithoutPatch]; ok {
				if value < patchCounter {
					data[imageWithoutPatch] = patchCounter
				}
				continue
			}

			data[imageWithoutPatch] = patchCounter
			continue
		}

		if _, ok := data[image]; ok {
			continue
		} else {
			data[image] = 0
		}

	}

	output := []string{}
	for k, v := range data {
		s := fmt.Sprintf("%s-%d", k, v)
		if v == 0 {
			output = append(output, k)
			continue
		}
		output = append(output, s)
	}
	return output
}
