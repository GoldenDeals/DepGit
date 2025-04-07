# Fix for "error: failed to push some refs to 'ssh://localhost:2222/test-repo.git'"

## Problem Identified

The error occurs because the `handleGitUploadPack` function in `internal/git/handler.go` was just a placeholder that returned an error:

```go
// handleGitUploadPack processes git fetch and clone operations
func handleGitUploadPack(s ssh.Session, repoName string) {
    // For now, this is a placeholder for future implementation
    io.WriteString(s, fmt.Sprintf("Repository '%s' fetch/clone not implemented yet\n", repoName))
    s.Exit(1)
}
```

When Git tries to push, it first needs to fetch the current state of the repository, which calls `git-upload-pack` on the server. Because this was not implemented, the push operation failed.

## Solution

The solution was to implement the `handleGitUploadPack` function properly, similar to how `handleGitReceivePack` is implemented:

```go
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

    // Initialize a bare git repository in the temp directory
    gitInitCmd := exec.Command("git", "init", "--bare", tempDir)
    if err := gitInitCmd.Run(); err != nil {
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

    // Execute git-upload-pack in the temp repository
    uploadPack := exec.Command("git", "upload-pack", tempDir)
    uploadPack.Stdin = s
    uploadPack.Stdout = s
    uploadPack.Stderr = s

    if err := uploadPack.Run(); err != nil {
        log.WithError(err).Error("git-upload-pack command failed")
        s.Exit(1)
        return
    }

    s.Exit(0)
}
```

Also, we had to update the call in the main handler:

```go
// Handle git-upload-pack (fetch/clone) command
if commandType == "git-upload-pack" {
    srv.handleGitUploadPack(s, repoName)
    return
}
```

## To Deploy

1. Rebuild and restart the server:
   ```
   cd cmd/app && go build && ./app
   ```

2. Try pushing to the repository again:
   ```
   bash scripts/test-git.sh
   ```

This implementation will properly handle the Git upload-pack protocol that is needed during push operations. 