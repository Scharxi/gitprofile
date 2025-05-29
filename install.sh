#!/bin/bash

# install.sh
set -e  # Exit on error

echo "🚀 Starting GitProfile Installation..."

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo "📦 Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
else
    echo "✅ Homebrew already installed"
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "📦 Installing Go..."
    brew install go
else
    echo "✅ Go already installed"
fi

echo "📦 Installing GitProfile..."
go install github.com/Scharxi/gitprofile@latest

echo "
🎉 Installation complete! 

To start using GitProfile:
1. Close and reopen your terminal or run: source ~/.zshrc
2. Create a new profile: gitprofile add work --name 'Your Name' --email 'your@email.com'
3. Use the profile in a repository: gitprofile use work

For more information, visit: https://github.com/yourusername/gitprofile
"