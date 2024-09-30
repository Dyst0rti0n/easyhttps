package easyhttps

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme"
)

// NewACMEClient creates a custom ACME client.
func NewACMEClient(directoryURL string, accountKey crypto.Signer) *acme.Client {
    if accountKey == nil {
        // Generate a new RSA private key for the account
        var err error
        accountKey, err = rsa.GenerateKey(rand.Reader, 2048)
        if err != nil {
            log.Fatalf("failed to generate account key: %v", err)
        }
    }

    return &acme.Client{
        Key:          accountKey,
        DirectoryURL: directoryURL,
        UserAgent:    "easyhttps",
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}
