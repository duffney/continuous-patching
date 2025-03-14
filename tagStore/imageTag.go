package tagStore

import (
	"fmt"
	"strconv"
	"strings"
)

type ImageTag struct {
	Name       string `json:"name"`
	BaseTag    string `json:"baseTag"`
	PatchLevel int    `json:"patchLevel"`
}

// NewImageTag parses a fully qualified image tag and creates an ImageTag struct
func NewImageTag(fullTag string) (*ImageTag, error) {
	parts := strings.Split(fullTag, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid image tag format: %s", fullTag)
	}

	name := parts[0]
	tag := parts[1]

	patchLevel := 0
	baseTag := tag

	if strings.Contains(tag, "-") {
		tagParts := strings.Split(tag, "-")
		if len(tagParts) != 2 {
			return nil, fmt.Errorf("invalid tag format: %s", tag)
		}

		baseTag = tagParts[0]
		patch, err := strconv.Atoi(tagParts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid patch number: %s", tagParts[1])
		}
		patchLevel = patch
	}

	return &ImageTag{
		Name:       name,
		BaseTag:    baseTag,
		PatchLevel: patchLevel,
	}, nil
}

func (it *ImageTag) GetCurrentTag() string {
	if it.PatchLevel == 0 {
		return fmt.Sprintf("%s:%s", it.Name, it.BaseTag)
	}
	return fmt.Sprintf("%s:%s-%d", it.Name, it.BaseTag, it.PatchLevel)
}

func (it *ImageTag) GetNextPatchTag() string {
	return fmt.Sprintf("%s:%s-%d", it.Name, it.BaseTag, it.PatchLevel+1)
}

func (it *ImageTag) IncrementPatch() {
	it.PatchLevel++
}
