#!/bin/bash

set -e

# Default values
APP_ADDRESS=${1:-"ssh://git@localhost:2222"}
REPO_NAME=${2:-"test-repo"}
PUSH_TO_SERVER=${3:-"yes"}  # Third parameter to control whether to push or not

TMP_DIR=$(mktemp -d -t "depgit-test-XXXXXX")

echo "Creating temporary git repository in: $TMP_DIR"
echo "Using DepGit server address: $APP_ADDRESS"
echo "Repository name: $REPO_NAME"
echo "Push to server: $PUSH_TO_SERVER"

# Function to clean up on exit
cleanup() {
  echo "Cleaning up temporary directory..."
  rm -rf "$TMP_DIR"
}

# Register cleanup function to be called on exit
trap cleanup EXIT

# Navigate to temporary directory
cd "$TMP_DIR"

# Initialize git repository
git init
git config --local commit.gpgsign false

# Create a README file with some content
cat > README.md << EOF
# Test Repository

This is a test repository for DepGit.

## About

This repository was automatically created by the test script.
EOF

# Configure git (required for commits)
git config user.email "test@example.com"
git config user.name "Test User"

# First commit
git add README.md
git commit -m "Initial commit"

# Add remote origin pointing to DepGit server
git remote add origin "$APP_ADDRESS/$REPO_NAME.git"

# Create additional files and make more commits
# Commit 2
echo "# File 1" > file1.txt
echo "This is the first additional file." >> file1.txt
git add file1.txt
git commit -m "Add file1.txt"

# Commit 3
mkdir -p src/app
cat > src/app/main.go << EOF
package main

import "fmt"

func main() {
    fmt.Println("Hello, DepGit!")
}
EOF
git add src/app/main.go
git commit -m "Add Go application"

# Commit 4
echo "*.log" > .gitignore
echo "tmp/" >> .gitignore
git add .gitignore
git commit -m "Add gitignore file"

# Commit 5
echo "# Development Notes" > NOTES.md
echo "* Remember to update documentation" >> NOTES.md
echo "* Fix all TODOs before release" >> NOTES.md
git add NOTES.md
git commit -m "Add development notes"

echo -e "\nRepository created with 5 commits."

# Check if we should push to the server
if [ "$PUSH_TO_SERVER" = "yes" ] || [ "$PUSH_TO_SERVER" = "y" ]; then
    echo -e "\nPushing to DepGit server at $APP_ADDRESS..."
    # Use -u to set upstream and 'master' is the branch name
    GIT_SSH_COMMAND="ssh" git push -u origin main
    echo -e "\nPush completed."
else
    echo -e "\nTo push to the DepGit server, run:"
    echo "cd $TMP_DIR && git push -u origin master"
fi

echo -e "\nRepository location: $TMP_DIR"
echo "Repository will be deleted when script exits. Copy to another location if you want to keep it."

# Keep the shell open with the repository until user presses a key
read -p "Press Enter to exit and clean up the repository..." 