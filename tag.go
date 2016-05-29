package libmason

import (
	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
)

// TagImage tags an image with newTag
func (h *DefaultHelper) TagImage(ctx context.Context, image string, newReference string) error {
	return h.client.ImageTag(ctx, image, newReference, types.ImageTagOptions{
		Force: true,
	})
}
