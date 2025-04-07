package git

import "github.com/gliderlabs/ssh"

type Config struct {
	Address string
}

func keyAuthOption(ctx ssh.Context, key ssh.PublicKey) bool {
	return true
}
