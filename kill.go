package libmason

import (
	"golang.org/x/net/context"
)

// ContainerKill stops the container execution abruptly.
func (h *DefaultHelper) ContainerKill(ctx context.Context, containerID string) error {
	return h.client.ContainerKill(context.Background(), containerID, "SIGKILL")
}
