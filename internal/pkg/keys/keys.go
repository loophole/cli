package keys

import (
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/terminal"
)

//ParsePublicKey retrieves an ssh.AuthMethod and the related PublicKey
func ParsePublicKey(terminalState *terminal.State, file string) (ssh.AuthMethod, ssh.PublicKey, error) { //key, err := ioutil.ReadFile(file)
	//https://godoc.org/golang.org/x/crypto/ssh/agent#ExtendedAgent
	// ssh-agent(1) provides a UNIX socket at $SSH_AUTH_SOCK.
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return nil, nil, err
	}

	agentClient := agent.NewClient(conn)

	signers, err := agentClient.Signers()

	signer := signers[0]

	return ssh.PublicKeys(signer), signer.PublicKey(), nil
}
