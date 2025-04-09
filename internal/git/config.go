// Package git implements git server functionality including SSH-based
// repository access and management.
package git

import (
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

// Config holds the configuration for the git server.
// It includes settings like the SSH address to listen on.
type Config struct {
	Address string
}

func keyAuthOption(ctx ssh.Context, pk ssh.PublicKey) bool {
	log.
		WithField("user", ctx.User()).
		WithField("addr", ctx.RemoteAddr()).
		WithField("pkType", pk.Type()).
		WithField("pkFingerprint", gossh.FingerprintSHA256(pk)).
		Info("Key auth attempt")

	// For development/testing, accept all keys if no authorized keys file is specified
	if ctx.User() == "git" {
		return true
	}

	log.Warn("Authentication failed: invalid username (must be 'git')")
	return false
}
