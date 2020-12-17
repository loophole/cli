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
	"github.com/loophole/cli/internal/pkg/urlmaker"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/webdav"
)

const (
	logoURL = "https://raw.githubusercontent.com/loophole/website/master/static/img/logo.png"
)

type ServerBuilder interface {
	WithHostname(string) ServerBuilder
	Proxy() ProxyServerBuilder
	ServeStatic() StaticServerBuilder
	ServeWebdav() WebdavServerBuilder
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
func (sb *serverBuilder) ServeWebdav() WebdavServerBuilder {
	return &webdavServerBuilder{
		serverBuilder: sb,
	}
}

// ProxyServerBuilder is used to proxy to already running server
type ProxyServerBuilder interface {
	ToEndpoint(lm.Endpoint) ProxyServerBuilder
	WithBasicAuth(string, string) ProxyServerBuilder
	DisableProxyErrorPage() ProxyServerBuilder
	Build() (*http.Server, error)
}
type proxyServerBuilder struct {
	serverBuilder         *serverBuilder
	endpoint              lm.Endpoint
	basicAuthEnabled      bool
	basicAuthUsername     string
	basicAuthPassword     string
	disableProxyErrorPage bool
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

func (psb *proxyServerBuilder) DisableProxyErrorPage() ProxyServerBuilder {
	psb.disableProxyErrorPage = true
	return psb
}

func (psb *proxyServerBuilder) Build() (*http.Server, error) {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: psb.endpoint.Protocol,
		Host:   psb.endpoint.Hostname(),
	})
	if !psb.disableProxyErrorPage {
		proxy.ErrorHandler = proxyErrorHandler
	}

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

// WebdavServerBuilder is used to create server which expose local directory
type WebdavServerBuilder interface {
	FromDirectory(string) WebdavServerBuilder
	WithBasicAuth(string, string) WebdavServerBuilder
	Build() (*http.Server, error)
}
type webdavServerBuilder struct {
	serverBuilder     *serverBuilder
	directory         string
	basicAuthEnabled  bool
	basicAuthUsername string
	basicAuthPassword string
}

func (wsb *webdavServerBuilder) FromDirectory(directory string) WebdavServerBuilder {
	wsb.directory = directory
	return wsb
}

func (wsb *webdavServerBuilder) WithBasicAuth(username string, password string) WebdavServerBuilder {
	wsb.basicAuthEnabled = true
	wsb.basicAuthUsername = username
	wsb.basicAuthPassword = password
	return wsb
}

func (wsb *webdavServerBuilder) Build() (*http.Server, error) {
	wdHandler := &webdav.Handler{
		Prefix:     "/",
		FileSystem: webdav.Dir(wsb.directory),
		LockSystem: webdav.NewMemLS(),
	}

	var server *http.Server

	if wsb.basicAuthEnabled {
		handler, err := getBasicAuthHandler(wsb.serverBuilder.siteID, wsb.basicAuthUsername, wsb.basicAuthPassword, wdHandler.ServeHTTP)
		if err != nil {
			return nil, err
		}

		server = &http.Server{
			Handler:   handler,
			TLSConfig: getTLSConfig(wsb.serverBuilder.siteID),
		}
	} else {
		server = &http.Server{
			Handler:   wdHandler,
			TLSConfig: getTLSConfig(wsb.serverBuilder.siteID),
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
		HostPolicy: autocert.HostWhitelist(urlmaker.GetSiteFQDN(siteID)),
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

	authenticator := auth.NewBasicAuthenticator(urlmaker.GetSiteFQDN(siteID), secret)
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

func proxyErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.Write([]byte(fmt.Sprintf(proxyErrorTemplate, logoURL, err.Error())))
}
