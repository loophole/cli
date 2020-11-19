package keys

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/terminal"
)

//ParsePublicKey retrieves an ssh.AuthMethod and the related PublicKey
func ParsePublicKey(terminalState *terminal.State, file string) (ssh.AuthMethod, ssh.PublicKey, error) {
	privateKey, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, nil, err
	}

	var passwordError *ssh.PassphraseMissingError
	var signer ssh.Signer
	signer, err = ssh.ParsePrivateKey(privateKey) //try to parse the key as if not password-protected

	if err != nil {
		if errors.As(err, &passwordError) { //if the key is password-protected, try to resolve it using the SSH-Agent, otherwise ask the user for the password
			publicKey, err := ioutil.ReadFile(file + ".pub")
			if err != nil {
				return nil, nil, err
			}

			signer, err = getSignerFromSSHAgent(publicKey)
			if err != nil {
				fmt.Print("Enter SSH password:")
				password, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
				fmt.Println()
				signer, err = ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(password))
				if err != nil {
					return nil, nil, err
				}
			}
		} else {
			return nil, nil, err
		}
	}

	return ssh.PublicKeys(signer), signer.PublicKey(), nil
}

//getSignerFromSSHAgent connects to the SSH Agent and tries to return a signer for the given publicKey
func getSignerFromSSHAgent(publicKey []byte) (ssh.Signer, error) {
	//https://godoc.org/golang.org/x/crypto/ssh/agent#ExtendedAgent
	// ssh-agent(1) provides a UNIX socket at $SSH_AUTH_SOCK.
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return nil, err
	}

	agentClient := agent.NewClient(conn)

	identities, err := agentClient.List()

	if err != nil {
		return nil, err
	}

	keyFound, index := keySavedInSSHAgent(publicKey, identities)

	if !keyFound {
		return nil, errors.New("key not found in SSH Agent")
	}

	signers, err := agentClient.Signers()

	if err != nil {
		return nil, err
	}

	return signers[index], nil
}

//keySavedInSSHAgent goes through the identities saved in SSH Agent and looks for a specific key. If found, it returns true and the index, otherwise false and -1.
func keySavedInSSHAgent(publicKey []byte, identities []*agent.Key) (result bool, index int) {
	for i, identity := range identities {
		publicKeyString := string(publicKey)[:len(string(publicKey))-1] //In my tests, string(publicKey) has 1 extra character added at the end which would make the comparison to identity.String() fail if not removed
		if publicKeyString == identity.String() {
			return true, i
		}
	}
	return false, -1
}
