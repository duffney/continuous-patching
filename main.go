// List all images and tags from GitHub Container Registry
// TODO: Add pagination
// TODO: Append registry name to image name
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("env:GITHUB_TOKEN is required")
	}

	userName := os.Getenv("GITHUB_USERNAME")
	if userName == "" {
		log.Fatal("env:GITHUB_USERNAME is required")
	}

	var url = "https://api.github.com/users/" + userName + "/packages?package_type=container"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var packages []Package
	err = json.Unmarshal(body, &packages)
	if err != nil {
		log.Fatal(err)
	}

	images := []string{}

	for _, p := range packages {
		url := fmt.Sprintf("https://api.github.com/users/%s/packages/%s/%s/versions", userName, p.PackageType, p.Name)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var tags []Tag
		err = json.Unmarshal(body, &tags)
		if err != nil {
			log.Fatal(err)
		}

		for _, t := range tags {
			for _, tag := range t.Metadata.Container.Tags {
				if !strings.Contains(tag, ".sig") {
					images = append(images, fmt.Sprintf("%s:%s", p.Name, tag))
				}
			}
		}
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
	fmt.Println(output)
	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	os.NewFile(0, "matrix.json")
	os.WriteFile("matrix.json", jsonOutput, 0644)
}
