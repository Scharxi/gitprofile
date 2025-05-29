package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestEnv(t *testing.T) (string, func()) {
	// Create a temporary directory for test config
	tmpDir, err := os.MkdirTemp("", "gitprofile-test-*")
	require.NoError(t, err)

	// Set the test config path
	testConfigPath := filepath.Join(tmpDir, ".gitprofiles.json")
	SetTestConfigPath(testConfigPath)

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tmpDir)
		SetTestConfigPath("") // Reset test config path
	}

	return tmpDir, cleanup
}

func getGitConfig(args ...string) (string, error) {
	output, err := runGitCommand(append([]string{"config", "--local"}, args...)...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func TestAddCommand(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	tests := []struct {
		name        string
		args        []string
		flags       map[string]string
		expectError bool
	}{
		{
			name: "valid profile",
			args: []string{"testprofile"},
			flags: map[string]string{
				"name":  "Test User",
				"email": "test@example.com",
			},
			expectError: false,
		},
		{
			name: "missing name",
			args: []string{"testprofile"},
			flags: map[string]string{
				"email": "test@example.com",
			},
			expectError: true,
		},
		{
			name: "missing email",
			args: []string{"testprofile"},
			flags: map[string]string{
				"name": "Test User",
			},
			expectError: true,
		},
		{
			name: "with gpg key",
			args: []string{"gpgprofile"},
			flags: map[string]string{
				"name":    "GPG User",
				"email":   "gpg@example.com",
				"gpg-key": "ABC123",
				"sign":    "true",
			},
			expectError: false,
		},
		{
			name: "with ssh key",
			args: []string{"sshprofile"},
			flags: map[string]string{
				"name":    "SSH User",
				"email":   "ssh@example.com",
				"ssh-key": "~/.ssh/id_rsa",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewAddCmd()
			for flag, value := range tt.flags {
				cmd.Flags().Set(flag, value)
			}

			err := cmd.RunE(cmd, tt.args)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Verify profile was saved
			profiles, err := LoadProfiles()
			require.NoError(t, err)

			profile, exists := profiles[tt.args[0]]
			assert.True(t, exists)
			assert.Equal(t, tt.flags["name"], profile.Name)
			assert.Equal(t, tt.flags["email"], profile.Email)
			if gpgKey, ok := tt.flags["gpg-key"]; ok {
				assert.Equal(t, gpgKey, profile.GPGKey)
			}
			if sshKey, ok := tt.flags["ssh-key"]; ok {
				assert.Equal(t, sshKey, profile.SSHKey)
			}
		})
	}
}

func TestListCommand(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	// Add some test profiles
	profiles := ProfileMap{
		"profile1": {
			Name:  "User 1",
			Email: "user1@example.com",
		},
		"profile2": {
			Name:        "User 2",
			Email:       "user2@example.com",
			GPGKey:      "ABC123",
			SignCommits: true,
		},
	}

	err := SaveProfiles(profiles)
	require.NoError(t, err)

	// Test list command
	cmd := NewListCmd()
	buffer := &bytes.Buffer{}
	cmd.SetOut(buffer)

	err = cmd.RunE(cmd, []string{})
	assert.NoError(t, err)

	output := buffer.String()
	assert.Contains(t, output, "profile1")
	assert.Contains(t, output, "profile2")
	assert.Contains(t, output, "user1@example.com")
	assert.Contains(t, output, "user2@example.com")
}

func TestUseCommand(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	// Create a test git repository
	repoDir := filepath.Join(tmpDir, "testrepo")
	err := os.MkdirAll(repoDir, 0755)
	require.NoError(t, err)

	// Initialize git repository
	err = os.Chdir(repoDir)
	require.NoError(t, err)
	_, err = runGitCommand("init")
	require.NoError(t, err)

	// Add a test profile
	profiles := ProfileMap{
		"testprofile": {
			Name:        "Test User",
			Email:       "test@example.com",
			GPGKey:      "ABC123",
			SignCommits: true,
			SSHKey:      "~/.ssh/id_rsa",
		},
	}

	err = SaveProfiles(profiles)
	require.NoError(t, err)

	// Test use command
	cmd := NewUseCmd()
	err = cmd.RunE(cmd, []string{"testprofile"})
	assert.NoError(t, err)

	// Verify git config
	name, err := getGitConfig("user.name")
	require.NoError(t, err)
	assert.Equal(t, "Test User", name)

	email, err := getGitConfig("user.email")
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", email)

	gpgKey, err := getGitConfig("user.signingkey")
	require.NoError(t, err)
	assert.Equal(t, "ABC123", gpgKey)

	signCommits, err := getGitConfig("commit.gpgsign")
	require.NoError(t, err)
	assert.Equal(t, "true", signCommits)

	sshCommand, err := getGitConfig("core.sshCommand")
	require.NoError(t, err)
	assert.Contains(t, sshCommand, "~/.ssh/id_rsa")
}

func TestStatusCommand(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	// Create a test git repository
	repoDir := filepath.Join(tmpDir, "testrepo")
	err := os.MkdirAll(repoDir, 0755)
	require.NoError(t, err)

	// Initialize git repository
	err = os.Chdir(repoDir)
	require.NoError(t, err)
	_, err = runGitCommand("init")
	require.NoError(t, err)

	// Add a test profile
	profiles := ProfileMap{
		"testprofile": {
			Name:  "Test User",
			Email: "test@example.com",
		},
	}

	err = SaveProfiles(profiles)
	require.NoError(t, err)

	// Use the profile
	useCmd := NewUseCmd()
	err = useCmd.RunE(useCmd, []string{"testprofile"})
	require.NoError(t, err)

	// Test status command
	cmd := NewStatusCmd()
	buffer := &bytes.Buffer{}
	cmd.SetOut(buffer)

	err = cmd.RunE(cmd, []string{})
	assert.NoError(t, err)

	output := buffer.String()
	assert.Contains(t, output, "Test User")
	assert.Contains(t, output, "test@example.com")
	assert.Contains(t, output, "testprofile")
}

func TestDeleteCommand(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	// Add test profiles
	profiles := ProfileMap{
		"profile1": {
			Name:  "User 1",
			Email: "user1@example.com",
		},
		"profile2": {
			Name:  "User 2",
			Email: "user2@example.com",
		},
	}

	err := SaveProfiles(profiles)
	require.NoError(t, err)

	tests := []struct {
		name        string
		profileName string
		expectError bool
	}{
		{
			name:        "delete existing profile",
			profileName: "profile1",
			expectError: false,
		},
		{
			name:        "delete non-existent profile",
			profileName: "nonexistent",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewDeleteCmd()
			err := cmd.RunE(cmd, []string{tt.profileName})

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Verify profile was deleted
			profiles, err := LoadProfiles()
			require.NoError(t, err)

			_, exists := profiles[tt.profileName]
			assert.False(t, exists)
		})
	}
}
