package libmason

import (
	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
)

// ContainerCreate creates a new Docker container and returns potential warnings
// FIXME(vdemeester) should validate options ?
func (h *DefaultHelper) ContainerCreate(ctx context.Context, createConfig types.ContainerCreateConfig) (types.ContainerCreateResponse, error) {
	return h.client.ContainerCreate(context.Background(), createConfig.Config, createConfig.HostConfig, createConfig.NetworkingConfig, createConfig.Name)
}
