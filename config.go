package easyhttps

import (
    "crypto/tls"
    "net/http"
    "time"

    "golang.org/x/crypto/acme"
    "golang.org/x/crypto/acme/autocert"
)

// Config holds configuration options for EasyHTTPS.
type Config struct {
    Domains            []string
    Email              string
    TLSConfig          *tls.Config
    TLSConfigCustomizer func(*tls.Config)
    HTTPHandler        http.Handler
    HTTPSAddr          string
    RedirectHTTP       bool
    ReadTimeout        time.Duration
    WriteTimeout       time.Duration
    ACMEClient         *acme.Client
    CertCache          autocert.Cache
    HostPolicy         autocert.HostPolicy
}

// defaultConfig provides default settings.
func defaultConfig() *Config {
    return &Config{
        HTTPSAddr:    ":443",
        RedirectHTTP: true,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        CertCache:    autocert.DirCache("certcache"),
    }
}

// Option is a function that configures Config.
type Option func(*Config)

// WithDomains sets the domains for which to obtain certificates.
func WithDomains(domains ...string) Option {
    return func(c *Config) {
        c.Domains = domains
    }
}

// WithEmail sets the email for Let's Encrypt account registration.
func WithEmail(email string) Option {
    return func(c *Config) {
        c.Email = email
    }
}

// WithTLSConfig allows custom TLS configuration.
func WithTLSConfig(tlsConfig *tls.Config) Option {
    return func(c *Config) {
        c.TLSConfig = tlsConfig
    }
}

// WithTLSConfigCustomizer allows fine-grained TLS settings.
func WithTLSConfigCustomizer(customizer func(*tls.Config)) Option {
    return func(c *Config) {
        c.TLSConfigCustomizer = customizer
    }
}

// WithHTTPHandler sets a custom HTTP handler for HTTP challenges.
func WithHTTPHandler(handler http.Handler) Option {
    return func(c *Config) {
        c.HTTPHandler = handler
    }
}

// WithHTTPSAddr sets the address for the HTTPS server.
func WithHTTPSAddr(addr string) Option {
    return func(c *Config) {
        c.HTTPSAddr = addr
    }
}

// WithRedirectHTTP sets whether to redirect HTTP to HTTPS.
func WithRedirectHTTP(redirect bool) Option {
    return func(c *Config) {
        c.RedirectHTTP = redirect
    }
}

// WithReadTimeout sets the server's read timeout.
func WithReadTimeout(timeout time.Duration) Option {
    return func(c *Config) {
        c.ReadTimeout = timeout
    }
}

// WithWriteTimeout sets the server's write timeout.
func WithWriteTimeout(timeout time.Duration) Option {
    return func(c *Config) {
        c.WriteTimeout = timeout
    }
}

// WithACMEClient allows setting a custom ACME client.
func WithACMEClient(client *acme.Client) Option {
    return func(c *Config) {
        c.ACMEClient = client
    }
}

// WithCertCache allows setting a custom certificate cache.
func WithCertCache(cache autocert.Cache) Option {
    return func(c *Config) {
        c.CertCache = cache
    }
}

// WithHostPolicy allows setting a custom host policy.
func WithHostPolicy(policy autocert.HostPolicy) Option {
    return func(c *Config) {
        c.HostPolicy = policy
    }
}

// newCertManager initializes the autocert.Manager with the given Config.
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
