package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Profile struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	GPGKey      string `json:"gpg_key,omitempty"`
	SignCommits bool   `json:"sign_commits"`
	SSHKey      string `json:"ssh_key,omitempty"`
}

type ProfileMap map[string]Profile

const configFileName = ".gitprofiles.json"

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFileName), nil
}

func LoadProfiles() (ProfileMap, error) {
	configPath, err := GetConfigPath()
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

func SaveProfiles(profiles ProfileMap) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetCurrentProfile returns the currently active profile in the current repository
func GetCurrentProfile() (*Profile, string, error) {
	// Get current git user name and email
	name, err := runGitCommand("config", "--local", "user.name")
	if err != nil {
		return nil, "", nil // No local git config found
	}

	email, err := runGitCommand("config", "--local", "user.email")
	if err != nil {
		return nil, "", nil
	}

	// Load all profiles
	profiles, err := LoadProfiles()
	if err != nil {
		return nil, "", err
	}

	// Find matching profile
	currentName := strings.TrimSpace(string(name))
	currentEmail := strings.TrimSpace(string(email))

	for profileName, profile := range profiles {
		if profile.Name == currentName && profile.Email == currentEmail {
			return &profile, profileName, nil
		}
	}

	return nil, "", nil
}
