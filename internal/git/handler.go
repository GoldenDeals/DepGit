package git

import (
	"io"

	"github.com/gliderlabs/ssh"
)

func handle(s ssh.Session) {
	_, err := io.WriteString(s, "HELLO USERNAME\n")
	if err != nil {
		log.WithError(err).Error("Failed to write to session")
		return
	}
}
