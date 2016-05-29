package libmason

import (
	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
)

// ContainerCommit creates a new Docker image from an existing Docker container.
func (h *DefaultHelper) ContainerCommit(ctx context.Context, container string, options types.ContainerCommitOptions) (string, error) {
	commitResponse, err := h.client.ContainerCommit(context.Background(), container, options)
	return commitResponse.ID, err
}
