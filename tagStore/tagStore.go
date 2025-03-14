package tagStore

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

type Matrix struct {
	Current_tag string `json:"current_tag"`
	Next_tag    string `json:"next_tag"`
}

// TagStore type that only uses a map
type TagStore struct {
	ImageMap map[string]*ImageTag // Map for storage and fast lookups
}

// NewImageRegistry creates a new registry
func New() *TagStore {
	return &TagStore{
		ImageMap: make(map[string]*ImageTag),
	}
}

// AddImage adds an image to the registry if it doesn't already exist
func (reg *TagStore) AddImage(fullTag string) error {
	img, err := NewImageTag(fullTag)
	if err != nil {
		return err
	}

	// Create a key for the map - using Name+BaseTag as the unique identifier
	key := img.Name + ":" + img.BaseTag

	// Only add if it doesn't exist
	if _, exists := reg.ImageMap[key]; !exists {
		reg.ImageMap[key] = img
	}

	return nil
}

// Flatten converts the map to a slice for JSON output
func (reg *TagStore) Flatten() []string {
	var images []string

	for _, img := range reg.ImageMap {
		images = append(images, img.GetCurrentTag())
	}

	return images
}

// TODO: include tag strategy as param
func (reg *TagStore) FlattenWithOptions(includeNextTag bool) interface{} {
	if includeNextTag {
		var m []Matrix
		for _, img := range reg.ImageMap {
			entry := Matrix{
				Current_tag: img.GetCurrentTag(),
				Next_tag:    img.GetNextPatchTag(),
			}
			m = append(m, entry)
		}
		return m
	} else {
		var images []string

		for _, img := range reg.ImageMap {
			images = append(images, img.GetCurrentTag())
		}

		return images
	}
}

// GetJSON returns the images as JSON
func (reg *TagStore) GetJSON(includeNextTag bool) (string, error) {
	data := reg.FlattenWithOptions(includeNextTag)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// GetImage retrieves an image if it exists
func (reg *TagStore) GetImage(name, baseTag string) (*ImageTag, bool) {
	key := name + ":" + baseTag
	img, exists := reg.ImageMap[key]
	return img, exists
}

// PatchImage increments the patch level for an image if it exists
func (reg *TagStore) PatchImage(name, baseTag string) (string, bool) {
	key := name + ":" + baseTag
	img, exists := reg.ImageMap[key]

	if !exists {
		return "", false
	}

	img.IncrementPatch()
	return img.GetCurrentTag(), true
}

func (reg *TagStore) PrintTable(includeNextTag bool) (string, error) {

	data := reg.FlattenWithOptions(includeNextTag)
	// First, convert the interface to JSON bytes
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Create table writer
	table := tablewriter.NewWriter(os.Stdout)

	// Try to unmarshal into slice of strings (first format)
	var simpleTags []string
	if err := json.Unmarshal(jsonBytes, &simpleTags); err == nil && len(simpleTags) > 0 {
		table.SetHeader([]string{"#", "Tag"})
		table.SetAutoFormatHeaders(false)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)

		for i, tag := range simpleTags {
			table.Append([]string{
				fmt.Sprintf("%d", i+1),
				tag,
			})
		}
		fmt.Println("Patchable Tags:")
		table.Render()
		return "", nil
	}

	// Try to unmarshal into slice of tag pairs (second format)
	type TagPair struct {
		CurrentTag string `json:"Current_tag"`
		NextTag    string `json:"Next_tag"`
	}
	var tagPairs []TagPair
	if err := json.Unmarshal(jsonBytes, &tagPairs); err == nil && len(tagPairs) > 0 {
		table.SetHeader([]string{"#", "Current Tag", "Next Tag"})
		table.SetAutoFormatHeaders(false)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)

		for i, pair := range tagPairs {
			table.Append([]string{
				fmt.Sprintf("%d", i+1),
				pair.CurrentTag,
				pair.NextTag,
			})
		}
		fmt.Println("Patchable Tags:")
		table.Render()
		return "", nil
	}

	// If neither format matches, just print the raw JSON
	fmt.Printf("Raw JSON:\n%s\n", string(jsonBytes))
	return "", nil
}
