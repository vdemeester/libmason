package libmason

import (
	"io"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/vdemeester/libmason/test"
)

type CopyClient struct {
	test.NopClient
	success bool
}

func (c *CopyClient) CopyToContainer(ctx context.Context, container, destPath string, content io.Reader, options types.CopyToContainerOptions) error {
	return c.NopClient.CopyToContainer(ctx, container, destPath, content, options)
}

func TestCopyToContainersErrors(t *testing.T) {
	client := &CopyClient{
		NopClient: test.NopClient{},
	}
	helper := &DefaultHelper{
		client: client,
	}
	err := helper.CopyToContainer(context.Background(), "container_id", "destpath", "srcpath", true)
	if err == nil {
		t.Fatalf("expected an error, got nothing")
	}
}
