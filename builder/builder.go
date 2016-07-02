package builder

import (
	"fmt"
	"sync"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types/strslice"
	"github.com/vdemeester/libmason"
)

// Config holds builder configuration (like the ImageID and additionnal attributes).
type Config struct {
	ImageID    string
	Entrypoint strslice.StrSlice
	Cmd        strslice.StrSlice

	attributes map[string]interface{}
	rwMu       sync.RWMutex
}

// Put adds a new attributes defined by a key and value.
func (c *Config) Put(key string, value interface{}) {
	c.rwMu.Lock()
	if c.attributes == nil {
		c.attributes = make(map[string]interface{})
	}
	c.attributes[key] = value
	c.rwMu.Unlock()
}

// Get returns the current attribute value and its existence from the specified key.
func (c *Config) Get(key string) (interface{}, bool) {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	value, exists := c.attributes[key]
	return value, exists
}

func noLog(_ string, _ ...interface{}) {}

// NewBuilder Creates a new step builder
func NewBuilder(helper libmason.Helper) *StepBuilder {
	return &StepBuilder{
		currentConfig: &Config{},
		helper:        helper,
		logfunc:       noLog,
	}
}

// WithSteps sets the steps to execute to a builder and returns it.
func WithSteps(builder *StepBuilder, steps []Step) *StepBuilder {
	builder.steps = steps
	return builder
}

// WithLogFunc sets the logging function to a builder and returns it.
// By default a builder has a noLog logging function (does nothing).
func WithLogFunc(builder *StepBuilder, fn func(string, ...interface{})) *StepBuilder {
	builder.logfunc = fn
	return builder
}

// StepBuilder is a builder that is composed of steps that are executed sequentially.
type StepBuilder struct {
	currentConfig *Config
	steps         []Step
	helper        libmason.Helper

	logfunc func(string, ...interface{})
}

// Run run the steps in order and returns the image ID generated.
// If a step fails, the run fails as well (at the first failure).
func (b *StepBuilder) Run(ctx context.Context) (string, error) {
	for stepNum, step := range b.steps {
		b.logfunc("Step %d: %s", stepNum, step)
		if step == nil {
			return "", fmt.Errorf("Step %d is nil", stepNum)
		}
		config, err := step.Execute(ctx, b.helper, b.currentConfig)
		if err != nil {
			return "", err
		}
		b.currentConfig = config
	}

	return b.currentConfig.ImageID, nil
}
