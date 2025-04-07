package git

import (
	"context"
	"os"

	"github.com/GoldenDeals/DepGit/internal/share/logger"
	"github.com/GoldenDeals/DepGit/internal/stroage"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

var log = logger.New("git")

type Server struct {
	config  Config
	srv     ssh.Server
	storage stroage.Storage
}

func Init(c Config, stroag stroage.Storage) (*Server, error) {
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

	s.srv.PublicKeyHandler = keyAuthOption
	s.srv.Handle(handle)

	return s, nil
}

func (s *Server) Serve(ctx context.Context) error {
	log.Debugf("Serving git ssh server on %s", s.config.Address)

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
