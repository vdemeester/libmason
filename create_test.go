package libmason

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"
	"github.com/vdemeester/libmason/test"
)

type CreateClient struct {
	test.NopClient
	response *types.ContainerCreateResponse
}

func (c *CreateClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (types.ContainerCreateResponse, error) {
	if c.response == nil {
		return c.NopClient.ContainerCreate(ctx, config, hostConfig, networkingConfig, containerName)
	}
	return *c.response, nil
}

func TestContainerCreateErrors(t *testing.T) {
	client := &CreateClient{
		NopClient: test.NopClient{},
	}
	helper := &DefaultHelper{
		client: client,
	}
	_, err := helper.ContainerCreate(context.Background(), types.ContainerCreateConfig{})
	if err == nil {
		t.Fatalf("expected an error, got nothing")
	}
}

func TestContainerCreate(t *testing.T) {
	client := &CreateClient{
		NopClient: test.NopClient{},
		response: &types.ContainerCreateResponse{
			ID: "container_id",
			Warnings: []string{
				"This is a warning",
			},
		},
	}
	helper := &DefaultHelper{
		client: client,
	}
	resp, err := helper.ContainerCreate(context.Background(), types.ContainerCreateConfig{})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "container_id" {
		t.Fatalf("expected response.ID to be 'container_id', got %q", resp.ID)
	}
	if len(resp.Warnings) != 1 {
		t.Fatalf("expected one warning, got %v", resp.Warnings)
	}
}
