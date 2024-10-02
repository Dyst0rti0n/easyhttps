# EasyHTTPS

[![GoDoc](https://godoc.org/github.com/Dyst0rti0n/easyhttps?status.svg)](https://godoc.org/github.com/Dyst0rti0n/easyhttps)
[![Go Report Card](https://goreportcard.com/badge/github.com/Dyst0rti0n/easyhttps)](https://goreportcard.com/report/github.com/Dyst0rti0n/easyhttps)
![Go Version](https://img.shields.io/github/go-mod/go-version/Dyst0rti0n/easyhttps)
![Contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

![Easyhttps Logo](./logo.png)

**EasyHTTPS** is a Go library that simplifies the process of setting up HTTPS servers with automatic certificate management using Let's Encrypt or any other ACME-compatible certificate authority. Whether you're working on small personal projects or large-scale applications, EasyHTTPS makes it easier than ever to enable HTTPS with minimal setup and advanced customization.

---

## Features ✨

- **Automatic HTTPS**: Instantly obtain and renew SSL/TLS certificates automatically.
- **Minimal Setup**: Start a secure server with as little as three lines of code.
- **Customisable**: Fine-grained control over TLS settings and certificate management.
- **ACME Compatibility**: Works with any ACME-compatible certificate authority, including Let's Encrypt.
- **Custom Certificate Storage**: Easily implement your own certificate storage backend.
- **Flexible Configuration**: Customise everything, from TLS settings to ACME client configurations.

---

## Table of Contents 📚

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Advanced Usage](#advanced-usage)
- [Configuration Options](#configuration-options)
- [Custom Certificate Storage](#custom-certificate-storage)
- [Integrating with Other ACME CAs](#integrating-with-other-acme-cas)
- [TLS Configuration](#tls-configuration)
- [License](#license)

---

## Installation ⚙️

To install **EasyHTTPS**, simply run the following command:

```bash
go get github.com/Dyst0rti0n/easyhttps
```

Or, if you're using Go modules:

```bash
go get github.com/Dyst0rti0n/easyhttps@latest
```

---

## Quick Start 🚀

Here's how to set up a basic HTTPS server with **minimal code**:

```go
package main

import (
    "net/http"

    "github.com/Dyst0rti0n/easyhttps"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, Secure World!"))
    })

    // Start the server with automatic HTTPS
    err := easyhttps.ListenAndServe(":80", nil)
    if err != nil {
        panic(err)
    }
}
```

This will:

- Serve your application on **port 443 (HTTPS)**.
- Automatically **obtain and renew certificates** from Let's Encrypt.
- **Redirect all HTTP traffic** on port 80 to HTTPS.

> **Note:** Ensure your domain's DNS records point to your server's IP address and that ports 80 and 443 are open.

---

## Advanced Usage ⚡

For more control over the server configuration, you can use various options:

```go
package main

import (
    "crypto/tls"
    "log"
    "net/http"

    "github.com/Dyst0rti0n/easyhttps"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, Advanced Secure World!"))
    })

    // Custom TLS configuration
    tlsCustomizer := func(tlsConfig *tls.Config) {
        tlsConfig.MinVersion = tls.VersionTLS13
        // Add additional TLS settings as needed
    }

    // Start the server with advanced options
    err := easyhttps.ListenAndServe(":80", mux,
        easyhttps.WithDomains("example.com", "www.example.com"),
        easyhttps.WithEmail("admin@example.com"),
        easyhttps.WithTLSConfigCustomizer(tlsCustomizer),
        easyhttps.WithRedirectHTTP(false),
        easyhttps.WithHTTPSAddr(":8443"),
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

---

## Configuration Options ⚙️

EasyHTTPS provides various options to customise its behavior:

- `WithDomains(domains ...string)`: Specify domain names for certificate issuance.
- `WithEmail(email string)`: Set the email address for ACME registration.
- `WithTLSConfig(tlsConfig *tls.Config)`: Provide a custom TLS configuration.
- `WithTLSConfigCustomizer(customizer func(*tls.Config))`: Customise TLS settings with a function.
- `WithHTTPHandler(handler http.Handler)`: Set a custom handler for HTTP challenges.
- `WithHTTPSAddr(addr string)`: Change the address for the HTTPS server (default `":443"`).
- `WithRedirectHTTP(redirect bool)`: Enable or disable HTTP to HTTPS redirection (default `true`).
- `WithReadTimeout(timeout time.Duration)`: Set the server's read timeout.
- `WithWriteTimeout(timeout time.Duration)`: Set the server's write timeout.
- `WithACMEClient(client *acme.Client)`: Use a custom ACME client for different CAs.
- `WithCertCache(cache autocert.Cache)`: Provide a custom certificate storage backend.
- `WithHostPolicy(policy autocert.HostPolicy)`: Set a custom host policy for domain validation.

---

## Custom Certificate Storage 🗄️

You can implement your own certificate storage backend by satisfying the `autocert.Cache` interface. This allows you to store certificates in a database, cloud storage, or any other storage solution.

Example:

```go
type CustomCache struct {
    // Implement your storage mechanism
}

func (cc *CustomCache) Get(ctx context.Context, name string) ([]byte, error) {
    // Retrieve data from storage
}

func (cc *CustomCache) Put(ctx context.Context, name string, data []byte) error {
    // Store data in storage
}

func (cc *CustomCache) Delete(ctx context.Context, name string) error {
    // Delete data from storage
}

// Usage
customCache := &CustomCache{}
err := easyhttps.ListenAndServe(":80", mux,
    easyhttps.WithCertCache(customCache),
)
```

---

## Integrating with Other ACME CAs 🔗

EasyHTTPS allows you to integrate with any ACME-compatible certificate authority by providing a custom ACME client:

```go
acmeClient := easyhttps.NewACMEClient("https://acme-staging-v02.api.letsencrypt.org/directory", nil)

err := easyhttps.ListenAndServe(":80", mux,
    easyhttps.WithACMEClient(acmeClient),
)
```

Replace the directory URL with the ACME server of your chosen certificate authority.

---

## TLS Configuration 🔐

For fine-grained control over TLS settings, use `WithTLSConfigCustomizer`:

```go
tlsCustomizer := func(tlsConfig *tls.Config) {
    tlsConfig.MinVersion = tls.VersionTLS13
    tlsConfig.CurvePreferences = []tls.CurveID{tls.X25519, tls.CurveP256}
    tlsConfig.CipherSuites = []uint16{
        tls.TLS_AES_128_GCM_SHA256,
        tls.TLS_AES_256_GCM_SHA384,
    }
}

err := easyhttps.ListenAndServe(":80", mux,
    easyhttps.WithTLSConfigCustomizer(tlsCustomizer),
)
```

---

## License 📝

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
