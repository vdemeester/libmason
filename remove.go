package libmason

import (
	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
)

// ContainerRemove removes a container specified by `id`.
func (h *DefaultHelper) ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error {
	return h.client.ContainerRemove(context.Background(), container, options)
}
