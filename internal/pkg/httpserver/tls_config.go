//go:build !dev
// +build !dev

package httpserver

import (
	"crypto/tls"
	"fmt"

	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/loophole/cli/internal/pkg/urlmaker"
	"golang.org/x/crypto/acme/autocert"
)

func getTLSConfig(siteID string, domain string, disableOldCiphers bool) *tls.Config {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(urlmaker.GetSiteFQDN(siteID, domain)),
		Cache:      autocert.DirCache(cache.GetLocalStorageDir("certs")),
		Email:      fmt.Sprintf("lh-%s@main.dev", siteID),
	}

	config := certManager.TLSConfig()
	if disableOldCiphers {
		config.MinVersion = tls.VersionTLS12
	}
	return config
}
