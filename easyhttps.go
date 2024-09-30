package easyhttps

import (
    "crypto/tls"
    "net/http"
    "strings"
)

// ListenAndServe starts the HTTPS server with automatic certificate management.
func ListenAndServe(addr string, handler http.Handler, options ...Option) error {
    cfg := defaultConfig()
    for _, opt := range options {
        opt(cfg)
    }

    // Initialize the certificate manager
    manager, err := cfg.newCertManager()
    if err != nil {
        return err
    }

    // Create TLS configuration
    tlsConfig := cfg.TLSConfig
    if tlsConfig == nil {
        tlsConfig = &tls.Config{
            GetCertificate: manager.GetCertificate,
            MinVersion:     tls.VersionTLS12,
            CipherSuites: []uint16{
                tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
                tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
                tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
                tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            },
        }
    }

    // Apply user-defined TLS settings
    if cfg.TLSConfigCustomizer != nil {
        cfg.TLSConfigCustomizer(tlsConfig)
    }

    httpsServer := &http.Server{
        Addr:         cfg.HTTPSAddr,
        Handler:      handler,
        TLSConfig:    tlsConfig,
        ReadTimeout:  cfg.ReadTimeout,
        WriteTimeout: cfg.WriteTimeout,
    }

    var httpServer *http.Server
    if cfg.RedirectHTTP {
        httpServer = &http.Server{
            Addr:         addr,
            Handler:      redirectHandler(),
            ReadTimeout:  cfg.ReadTimeout,
            WriteTimeout: cfg.WriteTimeout,
        }
    } else {
        httpServer = &http.Server{
            Addr:         addr,
            Handler:      manager.HTTPHandler(cfg.HTTPHandler),
            ReadTimeout:  cfg.ReadTimeout,
            WriteTimeout: cfg.WriteTimeout,
        }
    }

    errChan := make(chan error, 2)

    go func() {
        if err := httpServer.ListenAndServe(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
            errChan <- err
        }
    }()

    go func() {
        if err := httpsServer.ListenAndServeTLS("", ""); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
            errChan <- err
        }
    }()

    for i := 0; i < 2; i++ {
        if err := <-errChan; err != nil {
            return err
        }
    }

    return nil
}

func redirectHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        target := "https://" + r.Host + r.URL.RequestURI()
        http.Redirect(w, r, target, http.StatusMovedPermanently)
    })
}
