package libmason

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/vdemeester/libmason/test"
)

type KillClient struct {
	test.NopClient
	success bool
}

func (c *KillClient) ContainerKill(ctx context.Context, container, signal string) error {
	if c.success && signal == "SIGKILL" {
		return nil
	}
	return c.NopClient.ContainerKill(ctx, container, signal)
}

func TestContainerKillErrors(t *testing.T) {
	client := &KillClient{
		NopClient: test.NopClient{},
	}
	helper := &DefaultHelper{
		client: client,
	}
	err := helper.ContainerKill(context.Background(), "container_id")
	if err == nil {
		t.Fatalf("expected an error, got nothing")
	}
}

func TestContainerKill(t *testing.T) {
	client := &KillClient{
		NopClient: test.NopClient{},
		success:   true,
	}
	helper := &DefaultHelper{
		client: client,
	}
	err := helper.ContainerKill(context.Background(), "container_id")
	if err != nil {
		t.Fatal(err)
	}
}
