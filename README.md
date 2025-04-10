# GitProfile

A CLI tool to manage git profiles (name, email, GPG key, commit signing) and activate them per project.

## Installation

### Windows
1. Download the latest release from the releases page
2. Add the binary to your PATH

### macOS
```bash
brew install gitprofile
```

### Linux
```bash
# Download the binary
curl -L https://github.com/yourusername/gitprofile/releases/latest/download/gitprofile-linux-amd64 -o gitprofile
chmod +x gitprofile
sudo mv gitprofile /usr/local/bin/
```

### From Source
```bash
go install github.com/yourusername/gitprofile@latest
```

## Usage

### Add a new profile

```bash
gitprofile add work --name "John Doe" --email "john@company.com" --gpg-key "ABC123" --sign
```

### List all profiles

```bash
gitprofile list
```

### Use a profile in the current repository

```bash
gitprofile use work
```

## Features

- Store multiple git profiles with different configurations
- Set user name and email
- Configure GPG key and commit signing
- Apply profiles per repository
- Simple JSON-based storage in `~/.gitprofiles.json`
- Cross-platform support (Windows, macOS, Linux)

## Configuration File

The profiles are stored in `~/.gitprofiles.json` with the following structure:

```json
{
  "work": {
    "name": "John Doe",
    "email": "john@company.com",
    "gpg_key": "ABC123",
    "sign_commits": true
  },
  "personal": {
    "name": "John Doe",
    "email": "john@personal.com",
    "sign_commits": false
  }
}
```

## Platform Support

- Windows
- macOS
- Linux

## License

MIT 