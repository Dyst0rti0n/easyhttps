package easyhttps

import (
	"testing"
	"time"
)

func TestWithShutdownTimeout(t *testing.T) {
	initialTimeout := 5 * time.Second // Default value
	customTimeout := 10 * time.Second

	// Test default config
	cfg := defaultConfig()
	if cfg.ShutdownTimeout != initialTimeout {
		t.Errorf("defaultConfig ShutdownTimeout = %v, want %v", cfg.ShutdownTimeout, initialTimeout)
	}

	// Test WithShutdownTimeout option
	opt := WithShutdownTimeout(customTimeout)
	opt(cfg) // Apply the option
	if cfg.ShutdownTimeout != customTimeout {
		t.Errorf("WithShutdownTimeout failed to set ShutdownTimeout, got %v, want %v", cfg.ShutdownTimeout, customTimeout)
	}

	// Test applying to a fresh default config
	cfg2 := defaultConfig()
	options := []Option{WithShutdownTimeout(customTimeout)}
	for _, o := range options {
		o(cfg2)
	}
	if cfg2.ShutdownTimeout != customTimeout {
		t.Errorf("Applying WithShutdownTimeout via options slice failed, got %v, want %v", cfg2.ShutdownTimeout, customTimeout)
	}
}
