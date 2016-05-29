package libmason

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/vdemeester/libmason/test"
)

type StartClient struct {
	test.NopClient
	success bool
}

func (c *StartClient) ContainerStart(ctx context.Context, container, checkpoint string) error {
	if c.success {
		return nil
	}
	return c.NopClient.ContainerStart(ctx, container, checkpoint)
}

func TestContainerStartErrors(t *testing.T) {
	client := &StartClient{
		NopClient: test.NopClient{},
	}
	helper := &DefaultHelper{
		client: client,
	}
	err := helper.ContainerStart(context.Background(), "container_id")
	if err == nil {
		t.Fatalf("expected an error, got nothing")
	}
}

func TestContainerStart(t *testing.T) {
	client := &StartClient{
		NopClient: test.NopClient{},
		success:   true,
	}
	helper := &DefaultHelper{
		client: client,
	}
	err := helper.ContainerStart(context.Background(), "container_id")
	if err != nil {
		t.Fatal(err)
	}
}
