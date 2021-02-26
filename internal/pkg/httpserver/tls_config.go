// +build !dev

package httpserver

import (
	"crypto/tls"
	"fmt"

	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/loophole/cli/internal/pkg/urlmaker"
	"golang.org/x/crypto/acme/autocert"
)

func getTLSConfig(siteID string, domain string) *tls.Config {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(urlmaker.GetSiteFQDN(siteID, "loophole.site")),
		Cache:      autocert.DirCache(cache.GetLocalStorageDir("certs")),
		Email:      fmt.Sprintf("lh-%s@main.dev", siteID),
	}

	return certManager.TLSConfig()
}
