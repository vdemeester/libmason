package libmason

import (
	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
)

// ContainerStart starts a new container
func (h *DefaultHelper) ContainerStart(ctx context.Context, container string) error {
	return h.client.ContainerStart(context.Background(), container, types.ContainerStartOptions{})
}
