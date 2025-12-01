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

// GetFolderByName returns the folder with the given name, or nil if not found
func GetFolderByName(name string) *Folder {
	for _, f := range folders {
		if f.Name == name {
			return &f
		}
	}
	return nil
}

