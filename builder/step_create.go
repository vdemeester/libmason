package builder

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/vdemeester/libmason"
)

// ContainerID is the key used to store the current ContainerID in the builder attributes.
const ContainerID = "containerID"

// WithCreate creates a Create step with the specified step (and argumeents)
func WithCreate(step Step, entrypoint, cmd []string, stdin bool) Step {
	return &CreateStep{
		delegateStep: step,
		Entrypoint:   entrypoint,
		Cmd:          cmd,
		Stdin:        stdin,
	}
}

// CreateStep is a step that will create a container (based on the specified attributes) and
// and execute the specified delegate step in this container.
type CreateStep struct {
	delegateStep Step
	Cmd          []string
	Entrypoint   []string
	Stdin        bool
}

func (s *CreateStep) String() string {
	return fmt.Sprintf("%s (with create: %s %s)", s.delegateStep, s.Entrypoint, s.Cmd)
}

// Execute implements Step.Execute. It executes the step based on the specified config and helper.
func (s *CreateStep) Execute(ctx context.Context, helper libmason.Helper, config *Config) (*Config, error) {
	c, err := helper.ContainerCreate(ctx, types.ContainerCreateConfig{
		Config: &container.Config{
			Image:      config.ImageID,
			Entrypoint: s.Entrypoint,
			Cmd:        s.Cmd,
			OpenStdin:  s.Stdin,
			StdinOnce:  s.Stdin,
		},
	})
	if err != nil {
		return nil, err
	}

	config.Put(ContainerID, c.ID)
	return s.delegateStep.Execute(ctx, helper, config)
}
