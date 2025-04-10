package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Profile struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	GPGKey      string `json:"gpg_key,omitempty"`
	SignCommits bool   `json:"sign_commits"`
}

type ProfileMap map[string]Profile

const configFileName = ".gitprofiles.json"

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFileName), nil
}

func loadProfiles() (ProfileMap, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(ProfileMap), nil
		}
		return nil, err
	}

	var profiles ProfileMap
	if err := json.Unmarshal(data, &profiles); err != nil {
		return nil, err
	}

	return profiles, nil
}

func saveProfiles(profiles ProfileMap) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
