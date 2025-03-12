package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/duffney/copamatic/internal/ghcr"
	"github.com/spf13/cobra"
)

var outputFile string

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "",
	Long:  "",
	//TODO: Args: 2 (registryType ghcr|dockerhub|acr, name) ex copamatic registry ghcr duffney --list -o matrix.json
	// should these be flags instead?
	Run: handleRegistry,
}

func init() {
	registryCmd.Flags().BoolP("list", "l", false, "List patchable images from registry")
	registryCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
}

func handleRegistry(cmd *cobra.Command, args []string) {
	//TODO: Use docker token for registry auth (Next)
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("env:GITHUB_TOKEN is required")
	}

	userName := os.Getenv("GITHUB_USERNAME")
	if userName == "" {
		log.Fatal("env:GITHUB_USERNAME is required")
	}

	listFlag, _ := cmd.Flags().GetBool("list")
	if listFlag {
		patchTags := getPatchableTags(userName, token)

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
			//TODO: pretty print to table
			_, err := os.Stdout.Write([]byte(strings.Join(patchTags, "\n")))
			if err != nil {
				os.Exit(1)
			}
		}
	}

}

func getPatchableTags(username, token string) []string {
	//TODO: replace with interface
	images, err := ghcr.ListImages(username, token)
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
