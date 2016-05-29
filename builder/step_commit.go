package builder

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/vdemeester/libmason"
)

// WithCommit creates a Commit Step with the specified step.
func WithCommit(step Step) Step {
	return &CommitStep{
		delegateStep: step,
	}
}

// CommitStep is a step that execute the delegated step and the commit the current
// container into a new image.
type CommitStep struct {
	delegateStep Step
}

func (s *CommitStep) String() string {
	return fmt.Sprintf("%s (with commit)", s.delegateStep)
}

// Execute implements Step.Execute. It executes the step based on the specified config and helper.
func (s *CommitStep) Execute(ctx context.Context, helper libmason.Helper, config *Config) (*Config, error) {
	config, err := s.delegateStep.Execute(ctx, helper, config)
	if err != nil {
		return nil, err
	}

	containerID, ok := config.Get(ContainerID)
	if !ok {
		return nil, fmt.Errorf("%s missing in config, cannot commit the container", ContainerID)
	}
	imageID, err := helper.ContainerCommit(ctx, containerID.(string), types.ContainerCommitOptions{
		Changes: []string{
			fmt.Sprintf("CMD %v", config.Cmd),
			fmt.Sprintf("ENTRYPOINT %v", config.Entrypoint),
		},
	})
	if err != nil {
		return nil, err
	}

	config.ImageID = imageID

	return config, nil
}
