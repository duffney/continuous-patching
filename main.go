// List all images and tags from GitHub Container Registry
// TODO: Add pagination
// TODO: Add `latest` tag support for only the more recent image
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

//TODO: Add type to hold image and tags of hosted images
//TODO: Output images with all tag versions on separate lines

const url = "https://api.github.com/users/duffney/packages?package_type=container"

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is required")
	}

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

	// fmt.Println(string(body))
	var packages []Package
	err = json.Unmarshal(body, &packages)
	if err != nil {
		log.Fatal(err)
	}

	images := []string{}
	for _, p := range packages {
		url := fmt.Sprintf("https://api.github.com/users/duffney/packages/%s/%s/versions", p.PackageType, p.Name)
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
	output, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output))
	os.NewFile(0, "images.json")
	os.WriteFile("images.json", output, 0644)
}
