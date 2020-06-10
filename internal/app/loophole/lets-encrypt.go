package loophole

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"

	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/go-acme/lego/v3/certificate"
	"github.com/go-acme/lego/v3/challenge/http01"
	"github.com/go-acme/lego/v3/lego"
	"github.com/go-acme/lego/v3/registration"
)

func getLocalStorageDir() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Error reading user home directory %v", err)
	}

	return path.Join(home, ".local", "loophole")
}

// LEUser is type of a user that implements acme.User
type LEUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

// GetEmail is Email getter
func (u *LEUser) GetEmail() string {
	return u.Email
}

// GetRegistration is Registration getter
func (u LEUser) GetRegistration() *registration.Resource {
	return u.Registration
}

// GetPrivateKey is PrivateKey getter
func (u *LEUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// CheckCertificateExist checks if cert is already present
func CheckCertificateExist(siteID string) bool {
	storageDir := getLocalStorageDir()

	if fileExists(path.Join(storageDir, fmt.Sprintf("%s.crt", siteID))) && fileExists(path.Join(storageDir, fmt.Sprintf("%s.key", siteID))) {
		return true
	}
	return false
}

// GenerateCertificate generates Let's Encrypt certificate for domain
func GenerateCertificate(siteID string, operationStatusChannel chan bool) {
	// Create a user. New accounts need an email and private key to start.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	myUser := LEUser{
		Email: fmt.Sprintf("%s@main.dev", siteID),
		key:   privateKey,
	}

	config := lego.NewConfig(&myUser)

	// This CA URL is configured for a local dev instance of Boulder running in Docker in a VM.
	config.CADirURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
	config.Certificate.KeyType = certcrypto.RSA4096

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// We specify an http port of 5002 and an tls port of 5001 on all interfaces
	// because we aren't running as root and can't bind a listener to port 80 and 443
	// (used later when we attempt to pass challenges). Keep in mind that you still
	// need to proxy challenge traffic to port 5002 and 5001.
	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "5002"))
	if err != nil {
		log.Fatal(err)
	}
	// err = client.Challenge.SetTLSALPN01Provider(tlsalpn01.NewProviderServer("", "5001"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Fatal(err)
	}
	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{fmt.Sprintf("%s.loophole.site", siteID)},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Fatal(err)
	}

	// Each certificate comes back with the cert bytes, the bytes of the client's
	// private key, and a certificate URL. SAVE THESE TO DISK.
	fmt.Printf("%#v\n", certificates)

	close(operationStatusChannel)

	// ... all done.
}
