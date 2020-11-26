package keys

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/loophole/cli/internal/pkg/communication"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func ParsePublicKey(terminalState *terminal.State, file string) (ssh.AuthMethod, ssh.PublicKey, error) {
	key, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, nil, err
	}

	var passwordError *ssh.PassphraseMissingError
	signer, err := ssh.ParsePrivateKey(key)

	if err != nil {
		if errors.As(err, &passwordError) {
			communication.Write("Enter SSH password:")
			terminalState, err = terminal.GetState(int(os.Stdin.Fd()))
			if err != nil {
				return nil, nil, err
			}

			password, _ := terminal.ReadPassword(int(os.Stdin.Fd()))

			terminalState = nil

			fmt.Println()
			signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(password))
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	}

	return ssh.PublicKeys(signer), signer.PublicKey(), nil
}
