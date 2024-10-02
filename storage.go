package easyhttps

import (
    "context"
    "os"

    "golang.org/x/crypto/acme/autocert"
)

// An interface that extends autocert.Cache
type CustomCache interface {
    autocert.Cache
}

// Implements CustomCache using the filesystem
type FileCache struct {
    Dir string
}

func (fc FileCache) Get(ctx context.Context, name string) ([]byte, error) {
    return os.ReadFile(fc.filePath(name))
}

func (fc FileCache) Put(ctx context.Context, name string, data []byte) error {
    return os.WriteFile(fc.filePath(name), data, 0600)
}

func (fc FileCache) Delete(ctx context.Context, name string) error {
    return os.Remove(fc.filePath(name))
}

func (fc FileCache) filePath(name string) string {
    return fc.Dir + "/" + name
}

// Implements CustomCache in memory (not persistent)
type MemoryCache struct {
    m map[string][]byte
}

func NewMemoryCache() *MemoryCache {
    return &MemoryCache{
        m: make(map[string][]byte),
    }
}

func (mc *MemoryCache) Get(ctx context.Context, name string) ([]byte, error) {
    data, ok := mc.m[name]
    if !ok {
        return nil, os.ErrNotExist
    }
    return data, nil
}

func (mc *MemoryCache) Put(ctx context.Context, name string, data []byte) error {
    mc.m[name] = data
    return nil
}

func (mc *MemoryCache) Delete(ctx context.Context, name string) error {
    delete(mc.m, name)
    return nil
}
