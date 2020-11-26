package httpserver

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/cache"
	"golang.org/x/crypto/acme/autocert"
)

// NewProxy creates new reverse proxy for given endpoint
func NewProxy(localEndpoint lm.Endpoint, siteID string) *http.Server {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: localEndpoint.Protocol,
		Host:   localEndpoint.Hostname(),
	})

	server := &http.Server{
		Handler:   proxy,
		TLSConfig: getTLSConfig(siteID),
	}

	return server
}

// NewStaticServer creates new static server for given path
func NewStaticServer(localPath string, siteID string) *http.Server {
	fs := http.FileServer(http.Dir(localPath))

	server := &http.Server{
		Handler:   fs,
		TLSConfig: getTLSConfig(siteID),
	}

	return server
}

func getTLSConfig(siteID string) *tls.Config {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(fmt.Sprintf("%s.loophole.host", siteID)),
		Cache:      autocert.DirCache(cache.GetLocalStorageDir("certs")),
		Email:      fmt.Sprintf("lh-%s@main.dev", siteID),
	}

	return certManager.TLSConfig()
}
