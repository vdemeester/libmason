package builder

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/vdemeester/libmason"
)

// WithRemove creates a Remove Step with the specfied step.
func WithRemove(step Step) Step {
	return &RemoveStep{
		delegateStep: step,
	}
}

// RemoveStep is a step that execute the delegated step and remove the current container
type RemoveStep struct {
	delegateStep Step
}

func (s *RemoveStep) String() string {
	return fmt.Sprintf("%s (with remove)", s.delegateStep)
}

// Execute implements Step.Execute. It executes the step based on the specified config and helper.
func (s *RemoveStep) Execute(ctx context.Context, helper libmason.Helper, config *Config) (*Config, error) {
	config, err := s.delegateStep.Execute(ctx, helper, config)
	if err != nil {
		return nil, err
	}

	containerID, ok := config.Get(ContainerID)
	if !ok {
		return nil, fmt.Errorf("%s missing in config, cannot commit the container", ContainerID)
	}

	if err := helper.ContainerRemove(ctx, containerID.(string), types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		return nil, err
	}

	config.Put(ContainerID, "")
	return config, nil
}
