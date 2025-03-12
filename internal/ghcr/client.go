// TODO: Refactor into registry interface
package ghcr

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func ListImages(username, token string) ([]string, error) {
	var url = "https://api.github.com/users/" + username + "/packages?package_type=container"

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
		log.Fatalf("unexpected status code: %d, %s", resp.StatusCode, resp.Status)
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
		url := p.URL + "/versions"

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
			fmt.Println(url)
			fmt.Println(p.URL)
			fmt.Printf("package error, %s:%s:%s\n", p.Name, p.PackageType, p.URL)
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
	return images, nil
}
