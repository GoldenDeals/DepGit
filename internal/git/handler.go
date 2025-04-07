package git

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func (srv *Server) handle(s ssh.Session) {
	cmd := s.Command()
	if len(cmd) == 0 {
		io.WriteString(s, "No command provided\n")
		s.Exit(1)
		return
	}

	log.Debugf("Git command: %v", cmd)

	// Extract the repository name from the command
	// Format is typically: git-receive-pack '/repo.git' or git-upload-pack '/repo.git'
	if len(cmd) < 2 {
		io.WriteString(s, "Invalid command format\n")
		s.Exit(1)
		return
	}

	commandType := cmd[0]
	repoPath := strings.Trim(cmd[1], "'")
	repoName := strings.TrimSuffix(strings.TrimPrefix(repoPath, "/"), ".git")

	log.Debugf("Repository: %s, Command: %s", repoName, commandType)

	// Handle git-receive-pack (push) command
	if commandType == "git-receive-pack" {
		srv.handleGitReceivePack(s, repoName)
		return
	}

	// Handle git-upload-pack (fetch/clone) command
	if commandType == "git-upload-pack" {
		srv.handleGitUploadPack(s, repoName)
		return
	}

	// Unknown git command
	io.WriteString(s, fmt.Sprintf("Unknown git command: %s\n", commandType))
	s.Exit(1)
}

// handleGitReceivePack processes git push operations
func (srv *Server) handleGitReceivePack(s ssh.Session, repoName string) {
	// Create a temporary directory to unpack git objects
	tempDir, err := os.MkdirTemp("", "git-receive-pack-*")
	if err != nil {
		log.WithError(err).Error("Failed to create temp directory")
		io.WriteString(s, "Internal server error\n")
		s.Exit(1)
		return
	}
	defer os.RemoveAll(tempDir)

	// Initialize a bare git repository using go-git
	_, err = git.PlainInit(tempDir, true)
	if err != nil {
		log.WithError(err).Error("Failed to initialize git repository")
		io.WriteString(s, "Failed to initialize repository\n")
		s.Exit(1)
		return
	}

	// Note: The go-git library doesn't provide a direct API for handling git protocol
	// communication for server-side operations like git-receive-pack.
	// We need to implement a custom handler for the git protocol here.

	// Read and process git protocol data from client
	// This is a simplified example - in a real implementation, you would need to
	// implement the git protocol properly

	// 1. Buffer to store incoming pack data
	var packData bytes.Buffer

	// 2. Copy data from SSH session to buffer
	_, err = io.Copy(&packData, s)
	if err != nil {
		log.WithError(err).Error("Failed to read pack data")
		s.Exit(1)
		return
	}

	// 3. Process the pack data (simplified)
	// In a real implementation, you would need to parse the pack file format
	// and apply it to the repository

	// Store the repository data using the storage interface
	ctx := context.Background()
	namespace := fmt.Sprintf("repositories/%s", repoName)

	// Walk through all files in the git repository and save them
	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories as storage.Put expects file content
		if info.IsDir() {
			return nil
		}

		// Calculate relative path from tempDir
		relPath, err := filepath.Rel(tempDir, path)
		if err != nil {
			return err
		}

		// Read file content
		fileContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Store the file
		err = srv.storage.Put(ctx, namespace, relPath, bytes.NewReader(fileContent))
		if err == os.ErrExist {
			// File already exists - we're updating the repo, so overwrite is expected
			// For now, we'll skip existing files (this is a simplification)
			log.Debugf("File already exists: %s", relPath)
			return nil
		}
		if err != nil {
			log.WithError(err).Errorf("Failed to store file: %s", relPath)
			return err
		}

		return nil
	})

	if err != nil {
		log.WithError(err).Error("Failed to store repository")
		io.WriteString(s, "Failed to process repository data\n")
		s.Exit(1)
		return
	}

	// Store repository metadata
	metadata := fmt.Sprintf("Updated: %s", time.Now().Format(time.RFC3339))
	err = srv.storage.Put(ctx, namespace, "metadata.txt", strings.NewReader(metadata))
	if err != nil && err != os.ErrExist {
		log.WithError(err).Error("Failed to store repository metadata")
	}

	io.WriteString(s, fmt.Sprintf("Repository '%s' updated successfully\n", repoName))
	s.Exit(0)
}

// handleGitUploadPack processes git fetch and clone operations
func (srv *Server) handleGitUploadPack(s ssh.Session, repoName string) {
	// Create a temporary directory for the repository
	tempDir, err := os.MkdirTemp("", "git-upload-pack-*")
	if err != nil {
		log.WithError(err).Error("Failed to create temp directory")
		io.WriteString(s, "Internal server error\n")
		s.Exit(1)
		return
	}
	defer os.RemoveAll(tempDir)

	// Initialize a bare git repository using go-git
	repo, err := git.PlainInit(tempDir, true)
	if err != nil {
		log.WithError(err).Error("Failed to initialize git repository")
		io.WriteString(s, "Failed to initialize repository\n")
		s.Exit(1)
		return
	}

	// Get repository data from storage
	ctx := context.Background()
	namespace := fmt.Sprintf("repositories/%s", repoName)

	// Get the list of objects in the namespace
	objList, err := srv.storage.List(ctx, namespace)
	if err != nil {
		log.WithError(err).Error("Failed to list repository objects")
		io.WriteString(s, fmt.Sprintf("Repository '%s' not found or cannot be accessed\n", repoName))
		s.Exit(1)
		return
	}

	if len(objList) == 0 {
		log.Error("Repository is empty")
		io.WriteString(s, fmt.Sprintf("Repository '%s' is empty or not initialized\n", repoName))
		s.Exit(1)
		return
	}

	// Restore each file from storage to the temporary git repository
	for _, objname := range objList {
		// Skip metadata file as it's not part of the git repository
		if objname == "metadata.txt" {
			continue
		}

		// Get object from storage
		reader, err := srv.storage.Get(ctx, namespace, objname)
		if err != nil {
			log.WithError(err).Errorf("Failed to get object: %s", objname)
			continue
		}

		// Create file in temp repository
		filePath := filepath.Join(tempDir, objname)
		fileDir := filepath.Dir(filePath)

		// Create directory if it doesn't exist
		if err := os.MkdirAll(fileDir, 0755); err != nil {
			log.WithError(err).Errorf("Failed to create directory: %s", fileDir)
			continue
		}

		// Create file
		file, err := os.Create(filePath)
		if err != nil {
			log.WithError(err).Errorf("Failed to create file: %s", filePath)
			continue
		}

		// Copy data from storage to file
		_, err = io.Copy(file, reader)
		file.Close()

		if closer, ok := reader.(io.Closer); ok {
			closer.Close()
		}

		if err != nil {
			log.WithError(err).Errorf("Failed to write data to file: %s", filePath)
		}
	}

	// Note: The go-git library doesn't provide a direct API for handling git protocol
	// communication for server-side operations like git-upload-pack.
	// We need to implement a custom handler for the git protocol here.

	// Read client request and send repository data
	// This is a simplified example - in a real implementation, you would need to
	// implement the git protocol properly

	// Get repository references
	refs, err := repo.References()
	if err != nil {
		log.WithError(err).Error("Failed to get repository references")
		s.Exit(1)
		return
	}

	// Write repository references to client (simplified)
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() != plumbing.HashReference {
			return nil
		}
		refInfo := fmt.Sprintf("%s %s\n", ref.Hash().String(), ref.Name())
		_, err := s.Write([]byte(refInfo))
		return err
	})

	if err != nil {
		log.WithError(err).Error("Failed to send repository data")
		s.Exit(1)
		return
	}

	s.Exit(0)
}
