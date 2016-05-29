package builder

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/vdemeester/libmason"
)

// FromStep is the top-level step that should be. It's only valid if there is
// no other step executed before. It will get the specified image and put it
// into the builder config.
type FromStep struct {
	Reference string
}

func (s *FromStep) String() string {
	return fmt.Sprintf("FROM %s", s.Reference)
}

// Execute implements Step.Execute. It executes the step based on the specified config and helper.
func (s *FromStep) Execute(ctx context.Context, helper libmason.Helper, config *Config) (*Config, error) {
	// FromStep has to be the first step, so the config in argument should be nil
	if config.ImageID != "" {
		return nil, fmt.Errorf("From step should be the first step, was not.")
	}
	image, err := helper.GetImage(ctx, s.Reference, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	return &Config{
		ImageID: image.ID,
	}, nil
}
