package builder

import (
	"golang.org/x/net/context"

	"github.com/vdemeester/libmason"
)

// Step defines the method a builder step should define
type Step interface {
	// Execute the current step "content" using the specified helper and config.
	// The step should return the/an updated config (can be untouched too).
	Execute(ctx context.Context, helper libmason.Helper, config *Config) (*Config, error)
}

// NoopStep is a no-operation step, that does nothing
type NoopStep struct{}

// Execute implements Step.Execute. It executes the step based on the specified config and helper.
func (s *NoopStep) Execute(ctx context.Context, helper libmason.Helper, config *Config) (*Config, error) {
	return config, nil
}
