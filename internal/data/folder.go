package data

import (
	"encoding/json"
	"fmt"
	"os"
)

const foldersPath = "data/folders.json"

type Folder struct {
	Name string `json:"folder_name"`
	ID   string `json:"folder_id"`
}

var folders []Folder

// Load data from folders.json file
func Load() error {
	data, err := os.ReadFile(foldersPath)
	if err != nil {
		return fmt.Errorf("failed to read folders.json: %w", err)
	}
	if err := json.Unmarshal(data, &folders); err != nil {
		return fmt.Errorf("failed to parse folders.json: %w", err)
	}
	return nil
}

// GetFolders returns all folders
func GetFolders() []Folder {
	return folders
}

// GetFolderNames returns a slice of all folder names
func GetFolderNames() []string {
	names := make([]string, len(folders))
	for i, f := range folders {
		names[i] = f.Name
	}
	return names
}

// GetFolderIDsByNames returns the folder IDs with the given names
func GetFolderIDsByNames(names []string) ([]string, error) {
	ids := make([]string, 0)
	for _, name := range names {
		for _, f := range folders {
			if f.Name == name {
				ids = append(ids, f.ID)
				break
			}
		}
	}

	if len(ids) != len(names) {
		return nil, fmt.Errorf("some folders not found: %v", names)
	}

	return ids, nil
}
