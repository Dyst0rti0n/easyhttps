package easyhttps

import (
    "crypto/tls"
    "net/http"
    "time"

    "golang.org/x/crypto/acme"
    "golang.org/x/crypto/acme/autocert"
)

// Config options
type Config struct {
    Domains            []string
    Email              string
    TLSConfig          *tls.Config
    TLSConfigCustomiser func(*tls.Config)
    HTTPHandler        http.Handler
    HTTPSAddr          string
    RedirectHTTP       bool
    ReadTimeout        time.Duration
    WriteTimeout       time.Duration
    ACMEClient         *acme.Client
    CertCache          autocert.Cache
    HostPolicy         autocert.HostPolicy
}

// Default settings.
func defaultConfig() *Config {
    return &Config{
        HTTPSAddr:    ":443",
        RedirectHTTP: true,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        CertCache:    autocert.DirCache("certcache"),
    }
}

type Option func(*Config)

// Sets the domains for which to obtain certificates.
func WithDomains(domains ...string) Option {
    return func(c *Config) {
        c.Domains = domains
    }
}

// Sets the email for Let's Encrypt account registration.
func WithEmail(email string) Option {
    return func(c *Config) {
        c.Email = email
    }
}

// Allows custom TLS configuration.
func WithTLSConfig(tlsConfig *tls.Config) Option {
    return func(c *Config) {
        c.TLSConfig = tlsConfig
    }
}

// Allows fine-grained TLS settings.
func WithTLSConfigCustomiser(customiser func(*tls.Config)) Option {
    return func(c *Config) {
        c.TLSConfigCustomiser = customiser
    }
}

// sets a custom HTTP handler for HTTP challenges.
func WithHTTPHandler(handler http.Handler) Option {
    return func(c *Config) {
        c.HTTPHandler = handler
    }
}

// Sets the address for the HTTPS server.
func WithHTTPSAddr(addr string) Option {
    return func(c *Config) {
        c.HTTPSAddr = addr
    }
}

// Sets whether to redirect HTTP to HTTPS.
func WithRedirectHTTP(redirect bool) Option {
    return func(c *Config) {
        c.RedirectHTTP = redirect
    }
}

// Sets the server's read timeout.
func WithReadTimeout(timeout time.Duration) Option {
    return func(c *Config) {
        c.ReadTimeout = timeout
    }
}

// Sets the server's write timeout.
func WithWriteTimeout(timeout time.Duration) Option {
    return func(c *Config) {
        c.WriteTimeout = timeout
    }
}

// Allows setting a custom ACME client.
func WithACMEClient(client *acme.Client) Option {
    return func(c *Config) {
        c.ACMEClient = client
    }
}

// Allows setting a custom certificate cache.
func WithCertCache(cache autocert.Cache) Option {
    return func(c *Config) {
        c.CertCache = cache
    }
}

// Allows setting a custom host policy.
func WithHostPolicy(policy autocert.HostPolicy) Option {
    return func(c *Config) {
        c.HostPolicy = policy
    }
}

// Initialises the autocert.Manager with the given Config.
func (c *Config) newCertManager() (*autocert.Manager, error) {
    manager := &autocert.Manager{
        Prompt:     autocert.AcceptTOS,
        Cache:      c.CertCache,
        HostPolicy: c.HostPolicy,
        Email:      c.Email,
        Client:     c.ACMEClient,
    }

    // Set default HostPolicy if not provided
    if manager.HostPolicy == nil && len(c.Domains) > 0 {
        manager.HostPolicy = autocert.HostWhitelist(c.Domains...)
    }

    return manager, nil
}
