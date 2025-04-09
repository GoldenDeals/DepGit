package git

import (
	"context"
	"os"

	"github.com/GoldenDeals/DepGit/internal/share/errors"
	"github.com/GoldenDeals/DepGit/internal/share/logger"
	"github.com/GoldenDeals/DepGit/internal/stroage"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

var log = logger.New("git")

// Server represents a Git SSH server that handles repository operations.
// It manages SSH connections and delegates git commands to appropriate handlers.
type Server struct {
	config  Config
	srv     ssh.Server
	storage stroage.Storage
}

// Init creates and initializes a new Git SSH server with the given configuration
// and storage backend. It sets up SSH server settings and authentication.
func Init(c Config, stroag stroage.Storage) (*Server, error) {
	if c.Address == "" {
		log.Error("Address is empty")
		return nil, errors.ErrBadData
	}

	if stroag == nil {
		log.Error("Stroage is nil")
		return nil, errors.ErrBadData
	}
	s := new(Server)

	s.config = c
	s.storage = stroag
	s.srv = ssh.Server{
		Addr:   c.Address,
		Banner: "---------------- DepGit ----------------\n",
	}

	pemData, err := os.ReadFile(os.Getenv("DEPGIT_SSH_GIT_HOSTKEY"))
	if err != nil {
		return nil, err
	}

	sig, err := gossh.ParsePrivateKey(pemData)
	if err != nil {
		return nil, err
	}

	s.srv.AddHostKey(sig)

	//	s.srv.PublicKeyHandler = keyAuthOption
	s.srv.Handle(s.handler)

	return s, nil
}

// Serve starts the Git SSH server and listens for incoming connections.
// It blocks until the context is cancelled or an error occurs.
func (s *Server) Serve(ctx context.Context) error {
	log.
		WithField("addr", s.config.Address).
		Info("Serving git ssh server")

	errCh := make(chan error)
	go func() {
		errCh <- s.srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *Server) Close() error {
	log.Warn("Closing git ssh server")
	return s.srv.Close()
}
