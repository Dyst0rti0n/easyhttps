package easyhttps

import (
    "context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
    "time"
)

// Starts HTTPS server with auto cert management.
func ListenAndServe(addr string, handler http.Handler, options ...Option) error {
    cfg := defaultConfig()
    for _, opt := range options {
        opt(cfg)
    }

    // Initialise
    manager, err := cfg.newCertManager()
    if err != nil {
        return fmt.Errorf("failed to initialise cert manager: %w", err)
    }

    // Create TLS config
    tlsConfig := cfg.TLSConfig
    if tlsConfig == nil {
        tlsConfig = &tls.Config{
            GetCertificate: manager.GetCertificate,
            MinVersion:     tls.VersionTLS13,
            CipherSuites: []uint16{
                tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
                tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
                tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
                tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            },
        }
    }

    // Apply TLS settings
    if cfg.TLSConfigCustomiser != nil {
        cfg.TLSConfigCustomiser(tlsConfig)
    }

    httpsServer := &http.Server{
        Addr:         cfg.HTTPSAddr,
        Handler:      handler,
        TLSConfig:    tlsConfig,
        ReadTimeout:  cfg.ReadTimeout,
        WriteTimeout: cfg.WriteTimeout,
    }

    httpServer := &http.Server{
        Addr:         addr,
        Handler:      redirectHandler(),
        ReadTimeout:  cfg.ReadTimeout,
        WriteTimeout: cfg.WriteTimeout,
    }

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

    errChan := make(chan error, 2)
    go func() {
        errChan <- httpServer.ListenAndServe()
    }()

    go func() {
        errChan <- httpsServer.ListenAndServeTLS("", "")
    }()

    select {
    case err := <-errChan:
        return err
    case <-quit:
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        if err := httpServer.Shutdown(ctx); err != nil {
            return fmt.Errorf("HTTP server shutdown failed: %w", err)
        }
        if err := httpsServer.Shutdown(ctx); err != nil {
            return fmt.Errorf("HTTPS server shutdown failed: %w", err)
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
