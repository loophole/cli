package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/term"
)

//ParsePublicKey retrieves an ssh.AuthMethod and the related PublicKey
func ParsePublicKey(file string) (ssh.AuthMethod, ssh.PublicKey, error) {
	privateKey, err := ioutil.ReadFile(file)

	var pathError *os.PathError
	if errors.As(err, &pathError) { //if no keys are found, they are generated
		var publicKey []byte
		bitSize := 4096
		privateKey, publicKey, err = generateKeyPair(bitSize)
		if err != nil {
			return nil, nil, err
		}
		err := ioutil.WriteFile(file, privateKey, 0600)
		if err != nil {
			return nil, nil, err
		}

		err = ioutil.WriteFile(file+".pub", publicKey, 0600)
		if err != nil {
			return nil, nil, err
		}

	} else if err != nil {
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
				fmt.Print("Enter SSH password: ")

				password, _ := term.ReadPassword(int(os.Stdin.Fd()))
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

//adapted from https://gist.github.com/devinodaniel/8f9b8a4f31573f428f29ec0e884e6673
func generateKeyPair(bitSize int) (private []byte, public []byte, err error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, nil, err
	}

	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicKey)

	return privatePEM, pubKeyBytes, nil
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

//keySavedInSSHAgent goes through the identities saved in SSH Agent and looks for a specific key.
//If found, it returns true and the index, otherwise false and -1.
func keySavedInSSHAgent(publicKey []byte, identities []*agent.Key) (result bool, index int) {
	for i, identity := range identities {
		//In my tests, string(publicKey) has 1 extra character added at the end which would make the comparison to identity.String() fail if not removed
		publicKeyString := string(publicKey)[:len(string(publicKey))-1]
		if publicKeyString == identity.String() {
			return true, i
		}
	}
	return false, -1
}
