package httpserver

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	auth "github.com/abbot/go-http-auth"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/cache"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/bcrypt"
)

type ServerBuilder interface {
	WithHostname(string) ServerBuilder
	Proxy() ProxyServerBuilder
	ServeStatic() StaticServerBuilder
}

type serverBuilder struct {
	siteID string
}

func (sb *serverBuilder) WithHostname(siteID string) ServerBuilder {
	sb.siteID = siteID
	return sb
}

func (sb *serverBuilder) Proxy() ProxyServerBuilder {
	return &proxyServerBuilder{
		serverBuilder: sb,
	}
}

func (sb *serverBuilder) ServeStatic() StaticServerBuilder {
	return &staticServerBuilder{
		serverBuilder: sb,
	}
}

// ProxyServerBuilder is used to proxy to already running server
type ProxyServerBuilder interface {
	ToEndpoint(lm.Endpoint) ProxyServerBuilder
	WithBasicAuth(string, string) ProxyServerBuilder
	Build() (*http.Server, error)
}
type proxyServerBuilder struct {
	serverBuilder     *serverBuilder
	endpoint          lm.Endpoint
	basicAuthEnabled  bool
	basicAuthUsername string
	basicAuthPassword string
}

func (psb *proxyServerBuilder) ToEndpoint(endpoint lm.Endpoint) ProxyServerBuilder {
	psb.endpoint = endpoint
	return psb
}

func (psb *proxyServerBuilder) WithBasicAuth(username string, password string) ProxyServerBuilder {
	psb.basicAuthEnabled = true
	psb.basicAuthUsername = username
	psb.basicAuthPassword = password
	return psb
}

func (psb *proxyServerBuilder) Build() (*http.Server, error) {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: psb.endpoint.Protocol,
		Host:   psb.endpoint.Hostname(),
	})

	var server *http.Server

	if psb.basicAuthEnabled {
		proxyWithAuth, err := getBasicAuthHandler(psb.serverBuilder.siteID, psb.basicAuthUsername, psb.basicAuthPassword, proxy.ServeHTTP)
		if err != nil {
			return nil, err
		}

		server = &http.Server{
			Handler:   proxyWithAuth,
			TLSConfig: getTLSConfig(psb.serverBuilder.siteID),
		}
	} else {
		server = &http.Server{
			Handler:   proxy,
			TLSConfig: getTLSConfig(psb.serverBuilder.siteID),
		}
	}

	return server, nil
}

// StaticServerBuilder is used to create server which expose local directory
type StaticServerBuilder interface {
	FromDirectory(string) StaticServerBuilder
	WithBasicAuth(string, string) StaticServerBuilder
	Build() (*http.Server, error)
}
type staticServerBuilder struct {
	serverBuilder     *serverBuilder
	directory         string
	basicAuthEnabled  bool
	basicAuthUsername string
	basicAuthPassword string
}

func (ssb *staticServerBuilder) FromDirectory(directory string) StaticServerBuilder {
	ssb.directory = directory
	return ssb
}

func (ssb *staticServerBuilder) WithBasicAuth(username string, password string) StaticServerBuilder {
	ssb.basicAuthEnabled = true
	ssb.basicAuthUsername = username
	ssb.basicAuthPassword = password
	return ssb
}

func (ssb *staticServerBuilder) Build() (*http.Server, error) {
	fs := http.FileServer(http.Dir(ssb.directory))

	var server *http.Server

	if ssb.basicAuthEnabled {
		handler, err := getBasicAuthHandler(ssb.serverBuilder.siteID, ssb.basicAuthUsername, ssb.basicAuthPassword, fs.ServeHTTP)
		if err != nil {
			return nil, err
		}

		server = &http.Server{
			Handler:   handler,
			TLSConfig: getTLSConfig(ssb.serverBuilder.siteID),
		}
	} else {
		server = &http.Server{
			Handler:   fs,
			TLSConfig: getTLSConfig(ssb.serverBuilder.siteID),
		}
	}

	return server, nil
}

// New starts creation of new server
func New() ServerBuilder {
	return &serverBuilder{}
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

func getBasicAuthHandler(siteID string, username string, password string, handler http.HandlerFunc) (http.HandlerFunc, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	secret := getBasicAuthSecretParser(username, string(hashedPassword))

	authenticator := auth.NewBasicAuthenticator(fmt.Sprintf("%s.loophole.host", siteID), secret)
	return auth.JustCheck(authenticator, handler), nil
}

func getBasicAuthSecretParser(username string, hashedPassword string) auth.SecretProvider {
	return func(user string, realm string) string {
		if user == username {
			return hashedPassword
		}
		return ""
	}
}
