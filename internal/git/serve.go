package git

import (
	"context"
	"io"
	"os"

	"github.com/GoldenDeals/DepGit/internal/share/logger"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

var log = logger.New("git")

type Storage interface {
	Put(namespace string, objname string, obj io.Reader) error
	Get(namespace string, objname string) (io.Reader, error)
	List(namespace string) ([]string, error)
}

type Server struct {
	config Config
	srv    ssh.Server
	//stroage Storage
}

func Init(c Config) (*Server, error) {
	s := new(Server)

	s.config = c
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

	//s.stroage = stroage

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
