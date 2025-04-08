package git

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/gliderlabs/ssh"
)

func (s *Server) handler(conn ssh.Session) {
	// Get the git command from the SSH session
	cmd := conn.Command()
	if len(cmd) < 2 {
		log.
			WithContext(conn.Context()).
			WithField("command", cmd).
			WithField("user", conn.User()).
			WithField("addr", conn.RemoteAddr()).
			Debug("Empty command received")

		conn.Exit(1)
		return
	}

	gitCmd := cmd[0]

	if gitCmd == "git-receive-pack" {
		s.handleReceivePack(conn, cmd[1])
	} else if gitCmd == "git-upload-pack" {
		s.handleUploadPack(conn)
	} else {
		log.
			WithContext(conn.Context()).
			WithField("command", cmd).
			WithField("user", conn.User()).
			WithField("addr", conn.RemoteAddr()).
			Debug("Unknown git command received")

		conn.Exit(1)
		return
	}
}

func (s *Server) handleReceivePack(conn ssh.Session, repoName string) {
	log.
		WithContext(conn.Context()).
		WithField("user", conn.User()).
		WithField("addr", conn.RemoteAddr()).
		WithField("repo", repoName).
		Trace("Received git-receive-pack command")

	// Step 1: Advertise references - send list of refs to client
	// In a real implementation, this would get actual refs from the repository
	// For now, we'll just advertise a sample ref
	capabilities := "report-status delete-refs side-band-64k"
	ref := "refs/heads/main"
	objID := "0000000000000000000000000000000000000000" // Zero hash

	// Format: <len><objID> <ref>\0<capabilities>\n
	refLine := fmt.Sprintf("%s %s%c%s\n", objID, ref, 0, capabilities)
	pktLine := fmt.Sprintf("%04x%s", len(refLine)+4, refLine)

	// Write the reference advertisement
	conn.Write([]byte(pktLine))

	// Send flush packet to end the reference advertisement
	conn.Write([]byte("0000"))

	// Step 2: Read client commands
	reader := bufio.NewReader(conn)
	var commands []string

	// Read pkt-lines until flush packet (0000)
	for {
		// Read 4 bytes for the length
		prefix := make([]byte, 4)
		if _, err := io.ReadFull(reader, prefix); err != nil {
			log.
				WithContext(conn.Context()).
				WithField("user", conn.User()).
				WithField("addr", conn.RemoteAddr()).
				WithField("repo", repoName).
				WithError(err).
				Debug("Failed to read packet prefix")
			conn.Exit(1)
			return
		}

		// Convert hex length to int
		length := 0
		if _, err := fmt.Sscanf(string(prefix), "%04x", &length); err != nil {
			log.
				WithContext(conn.Context()).
				WithField("user", conn.User()).
				WithField("addr", conn.RemoteAddr()).
				WithField("repo", repoName).
				WithError(err).
				Debug("Failed to parse packet length")
			conn.Exit(1)
			return
		}

		// Check for flush packet
		if length == 0 {
			break
		}

		// Read the actual data (length includes the 4 bytes we already read)
		data := make([]byte, length-4)
		if _, err := io.ReadFull(reader, data); err != nil {
			log.
				WithContext(conn.Context()).
				WithField("user", conn.User()).
				WithField("addr", conn.RemoteAddr()).
				WithField("repo", repoName).
				WithError(err).
				Debug("Failed to read packet data")
			conn.Exit(1)
			return
		}

		commands = append(commands, string(data))
		log.
			WithContext(conn.Context()).
			WithField("command", string(data)).
			WithField("repo", repoName).
			Trace("Received command")
	}

	// Step 3: Parse packfile
	log.
		WithContext(conn.Context()).
		WithField("user", conn.User()).
		WithField("addr", conn.RemoteAddr()).
		WithField("repo", repoName).
		Debug("Parsing packfile")

	objs, err := ParsePackfile(reader)
	if err != nil {
		log.
			WithContext(conn.Context()).
			WithField("user", conn.User()).
			WithField("addr", conn.RemoteAddr()).
			WithField("repo", repoName).
			WithError(err).
			Error("Failed to parse packfile")

		// Send error response
		errorMsg := fmt.Sprintf("unpack error: %s", err)
		conn.Write([]byte(fmt.Sprintf("%04x%s", len(errorMsg)+4, errorMsg)))
		conn.Write([]byte("0000"))
		conn.Exit(1)
		return
	}

	log.
		WithContext(conn.Context()).
		WithField("user", conn.User()).
		WithField("addr", conn.RemoteAddr()).
		WithField("repo", repoName).
		WithField("objectCount", len(objs)).
		Info("Successfully parsed packfile")

	// Log parsed objects
	for _, obj := range objs {
		log.
			WithContext(conn.Context()).
			WithField("user", conn.User()).
			WithField("addr", conn.RemoteAddr()).
			WithField("repo", repoName).
			WithField("objHash", obj.Hash()).
			WithField("objType", obj.Type()).
			WithField("objSize", obj.Size()).
			Debug("Parsed object")
	}

	// Step 4: Send success response
	// In real implementation, we would process the commands here
	// Format responses according to pkt-line

	// Unpack OK message
	unpackOk := "unpack ok\n"
	conn.Write([]byte(fmt.Sprintf("%04x%s", len(unpackOk)+4, unpackOk)))

	// Command status (for each reference update)
	for _, cmd := range commands {
		parts := strings.Fields(cmd)
		if len(parts) >= 2 {
			refName := parts[1]
			status := fmt.Sprintf("ok %s\n", refName)
			conn.Write([]byte(fmt.Sprintf("%04x%s", len(status)+4, status)))
		}
	}

	// Send flush packet to end the response
	conn.Write([]byte("0000"))

	log.
		WithContext(conn.Context()).
		WithField("user", conn.User()).
		WithField("addr", conn.RemoteAddr()).
		WithField("repo", repoName).
		Info("Completed receive-pack successfully")

	conn.Exit(0)
}

func (s *Server) handleUploadPack(conn ssh.Session) {
	log.
		WithContext(conn.Context()).
		WithField("user", conn.User()).
		WithField("addr", conn.RemoteAddr()).
		Trace("Received git-upload-pack command")

	panic("not implemented")
}
