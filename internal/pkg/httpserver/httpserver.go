package httpserver

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	auth "github.com/abbot/go-http-auth"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/urlmaker"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/webdav"
)

const (
	logoURL = "https://raw.githubusercontent.com/loophole/website/master/static/img/logo.png"
)

type ServerBuilder interface {
	WithSiteID(string) ServerBuilder
	WithDomain(string) ServerBuilder
	DisableOldCiphers(bool) ServerBuilder
	Proxy() ProxyServerBuilder
	ServeStatic() StaticServerBuilder
	ServeWebdav() WebdavServerBuilder
}

type serverBuilder struct {
	siteID            string
	domain            string
	disableOldCiphers bool
}

func (sb *serverBuilder) WithSiteID(siteID string) ServerBuilder {
	sb.siteID = siteID
	return sb
}

func (sb *serverBuilder) WithDomain(domain string) ServerBuilder {
	sb.domain = domain
	return sb
}
func (sb *serverBuilder) DisableOldCiphers(setting bool) ServerBuilder {
	sb.disableOldCiphers = setting
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
	EnableInsecureHTTPSBackend() ProxyServerBuilder
	Build() (*http.Server, error)
}
type proxyServerBuilder struct {
	serverBuilder         *serverBuilder
	endpoint              lm.Endpoint
	basicAuthEnabled      bool
	basicAuthUsername     string
	basicAuthPassword     string
	disableProxyErrorPage bool
	disableCertCheck      bool
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

func (psb *proxyServerBuilder) EnableInsecureHTTPSBackend() ProxyServerBuilder {
	psb.disableCertCheck = true
	return psb
}

func (psb *proxyServerBuilder) Build() (*http.Server, error) {
	target := &url.URL{
		Scheme: psb.endpoint.Protocol,
		Host:   psb.endpoint.Hostname(),
	}
	if psb.endpoint.Path != "" {
		target.Path = psb.endpoint.Path
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	defaultDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		defaultDirector(req)

		addr := net.ParseIP(target.Host)
		if addr == nil {
			req.Host = target.Host
		}

		req.Header.Set("X-Forwarded-Host", urlmaker.GetSiteFQDN(psb.serverBuilder.siteID, psb.serverBuilder.domain))
		req.Header.Set("X-Forwarded-Proto", "https")
	}

	if !psb.disableProxyErrorPage {
		proxy.ErrorHandler = proxyErrorHandler
	}

	if psb.disableCertCheck {
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	var server *http.Server

	if psb.basicAuthEnabled {
		proxyWithAuth, err := getBasicAuthHandler(psb.serverBuilder.siteID, psb.serverBuilder.domain, psb.basicAuthUsername, psb.basicAuthPassword, proxy.ServeHTTP)
		if err != nil {
			return nil, err
		}

		server = &http.Server{
			Handler:   proxyWithAuth,
			TLSConfig: getTLSConfig(psb.serverBuilder.siteID, psb.serverBuilder.domain, psb.serverBuilder.disableOldCiphers),
		}
	} else {
		server = &http.Server{
			Handler:   proxy,
			TLSConfig: getTLSConfig(psb.serverBuilder.siteID, psb.serverBuilder.domain, psb.serverBuilder.disableOldCiphers),
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
		handler, err := getBasicAuthHandler(ssb.serverBuilder.siteID, ssb.serverBuilder.domain, ssb.basicAuthUsername, ssb.basicAuthPassword, fs.ServeHTTP)
		if err != nil {
			return nil, err
		}

		server = &http.Server{
			Handler:   handler,
			TLSConfig: getTLSConfig(ssb.serverBuilder.siteID, ssb.serverBuilder.domain, ssb.serverBuilder.disableOldCiphers),
		}
	} else {
		server = &http.Server{
			Handler:   fs,
			TLSConfig: getTLSConfig(ssb.serverBuilder.siteID, ssb.serverBuilder.domain, ssb.serverBuilder.disableOldCiphers),
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
		handler, err := getBasicAuthHandler(wsb.serverBuilder.siteID, wsb.serverBuilder.domain, wsb.basicAuthUsername, wsb.basicAuthPassword, wdHandler.ServeHTTP)
		if err != nil {
			return nil, err
		}

		server = &http.Server{
			Handler:   handler,
			TLSConfig: getTLSConfig(wsb.serverBuilder.siteID, wsb.serverBuilder.domain, wsb.serverBuilder.disableOldCiphers),
		}
	} else {
		server = &http.Server{
			Handler:   wdHandler,
			TLSConfig: getTLSConfig(wsb.serverBuilder.siteID, wsb.serverBuilder.domain, wsb.serverBuilder.disableOldCiphers),
		}
	}

	return server, nil
}

// New starts creation of new server
func New() ServerBuilder {
	return &serverBuilder{}
}

func getBasicAuthHandler(siteID string, domain string, username string, password string, handler http.HandlerFunc) (http.HandlerFunc, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	secret := getBasicAuthSecretParser(username, string(hashedPassword))

	authenticator := auth.NewBasicAuthenticator(urlmaker.GetSiteFQDN(siteID, domain), secret)
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
