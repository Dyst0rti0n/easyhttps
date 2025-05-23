package easyhttps

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "fmt"
    "log"
    "net/http"
    "time"

    "golang.org/x/crypto/acme"
)

// Creates a custom ACME client.
func NewACMEClient(directoryURL string, accountKey crypto.Signer) (*acme.Client, error) {
    if accountKey == nil {
        // Generate a new RSA PK
        var err error
        accountKey, err = rsa.GenerateKey(rand.Reader, 2048)
        if err != nil {
            return nil, fmt.Errorf("failed to generate account key: %w", err)
        }
    }

    return &acme.Client{
        Key:          accountKey,
        DirectoryURL: directoryURL,
        UserAgent:    "easyhttps",
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }, nil
}
