//go:build dev
// +build dev

package httpserver

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"strings"
	"time"

	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/pkg/errors"
)

func getTLSConfig(siteID string, domain string, disableOldCiphers bool) *tls.Config {
	config := &tls.Config{
		GetCertificate: getCertificate(fmt.Sprintf("%s.%s", siteID, domain)),
	}
	if disableOldCiphers {
		config.MinVersion = tls.VersionTLS12
	}
	return config
}

func getCertificate(arg interface{}) func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	var opts certopts
	var err error
	if host, ok := arg.(string); ok {
		opts = certopts{
			RsaBits:   2048,
			Host:      host,
			ValidFrom: time.Now(),
		}
	} else if o, ok := arg.(certopts); ok {
		opts = o
	} else {
		err = errors.New("Invalid arg type, must be string(hostname) or Certopt{...}")
	}
	cert, err := generate(opts)
	return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		if err != nil {
			return nil, err
		}
		communication.Info("Obtained development certificate")
		return cert, nil
	}
}

type certopts struct {
	RsaBits   int
	Host      string
	IsCA      bool
	ValidFrom time.Time
	ValidFor  time.Duration
}

func generate(opts certopts) (*tls.Certificate, error) {

	priv, err := rsa.GenerateKey(rand.Reader, opts.RsaBits)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate private key")
	}

	notAfter := opts.ValidFrom.Add(opts.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate serial number\n")
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: opts.ValidFrom,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(opts.Host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if opts.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create certificate")
	}

	return &tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  priv,
	}, nil
}
