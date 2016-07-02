package libmason

import (
	"golang.org/x/net/context"
)

// TagImage tags an image with newTag
func (h *DefaultHelper) TagImage(ctx context.Context, image string, newReference string) error {
	return h.client.ImageTag(ctx, image, newReference)
}
