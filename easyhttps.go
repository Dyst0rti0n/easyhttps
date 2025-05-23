package easyhttps

import (
    "context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

    runErrChan := make(chan error, 2) // Renamed and still buffered
    go func() {
        if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            runErrChan <- fmt.Errorf("HTTP server ListenAndServe error: %w", err)
        }
    }()

    go func() {
        if err := httpsServer.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
            runErrChan <- fmt.Errorf("HTTPS server ListenAndServeTLS error: %w", err)
        }
    }()

    var startupError error
    select {
    case err := <-runErrChan:
        startupError = err
        close(quit) // Trigger shutdown path
    case <-quit:
        // OS signal received, proceed to shutdown
    }

    // Shutdown logic
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var httpShutdownErr, httpsShutdownErr error

    if err := httpServer.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
        httpShutdownErr = err
    }
    if err := httpsServer.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
        httpsShutdownErr = err
    }

    // Consolidate errors for return
    var errorMessages []string
    if startupError != nil {
        errorMessages = append(errorMessages, startupError.Error())
    }
    if httpShutdownErr != nil {
        errorMessages = append(errorMessages, fmt.Sprintf("HTTP server shutdown failed: %v", httpShutdownErr))
    }
    if httpsShutdownErr != nil {
        errorMessages = append(errorMessages, fmt.Sprintf("HTTPS server shutdown failed: %v", httpsShutdownErr))
    }

    if len(errorMessages) > 0 {
        return fmt.Errorf(strings.Join(errorMessages, "; "))
    }

    return nil
}

func redirectHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        target := "https://" + r.Host + r.URL.RequestURI()
        http.Redirect(w, r, target, http.StatusMovedPermanently)
    })
}
