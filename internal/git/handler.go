package git

import (
	"io"

	"github.com/gliderlabs/ssh"
)

func handle(s ssh.Session) {
	io.WriteString(s, "HELLO USERNAME\n")
}
