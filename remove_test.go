package libmason

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/vdemeester/libmason/test"
)

type RemoveClient struct {
	test.NopClient
	success bool
}

func (c *RemoveClient) ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error {
	if c.success {
		return nil
	}
	return c.NopClient.ContainerRemove(ctx, container, options)
}

func TestContainerRemoveErrors(t *testing.T) {
	client := &RemoveClient{
		NopClient: test.NopClient{},
	}
	helper := &DefaultHelper{
		client: client,
	}
	err := helper.ContainerRemove(context.Background(), "container_id", types.ContainerRemoveOptions{})
	if err == nil {
		t.Fatalf("expected an error, got nothing")
	}
}

func TestContainerRemove(t *testing.T) {
	client := &RemoveClient{
		NopClient: test.NopClient{},
		success:   true,
	}
	helper := &DefaultHelper{
		client: client,
	}
	err := helper.ContainerRemove(context.Background(), "container_id", types.ContainerRemoveOptions{})
	if err != nil {
		t.Fatal(err)
	}
}
