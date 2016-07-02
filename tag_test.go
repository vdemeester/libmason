package libmason

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/vdemeester/libmason/test"
)

type TagClient struct {
	test.NopClient
	success bool
}

func (c *TagClient) ImageTag(ctx context.Context, image, newReference string) error {
	if c.success {
		return nil
	}
	return c.NopClient.ImageTag(ctx, image, newReference)
}

func TestTagImageErrors(t *testing.T) {
	client := &TagClient{
		NopClient: test.NopClient{},
	}
	helper := &DefaultHelper{
		client: client,
	}
	err := helper.TagImage(context.Background(), "image_id", "reference")
	if err == nil {
		t.Fatalf("expected an error, got nothing")
	}
}

func TestTagImage(t *testing.T) {
	client := &TagClient{
		NopClient: test.NopClient{},
		success:   true,
	}
	helper := &DefaultHelper{
		client: client,
	}
	err := helper.TagImage(context.Background(), "image_id", "reference")
	if err != nil {
		t.Fatal(err)
	}
}
