package libmason

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/vdemeester/libmason/test"
)

type CommitClient struct {
	test.NopClient
	success bool
}

func (c *CommitClient) ContainerCommit(ctx context.Context, container string, options types.ContainerCommitOptions) (types.ContainerCommitResponse, error) {
	if c.success {
		return types.ContainerCommitResponse{
			ID: "ID",
		}, nil
	}
	return c.NopClient.ContainerCommit(ctx, container, options)
}

func TestContainerCommitErrors(t *testing.T) {
	client := &CommitClient{
		NopClient: test.NopClient{},
	}
	helper := &DefaultHelper{
		client: client,
	}
	_, err := helper.ContainerCommit(context.Background(), "container_id", types.ContainerCommitOptions{})
	if err == nil {
		t.Fatalf("expected an error, got nothing")
	}
}

func TestContainerCommit(t *testing.T) {
	client := &CommitClient{
		NopClient: test.NopClient{},
		success:   true,
	}
	helper := &DefaultHelper{
		client: client,
	}
	id, err := helper.ContainerCommit(context.Background(), "container_id", types.ContainerCommitOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if id != "ID" {
		t.Fatalf("expected id ID, got %s", id)
	}
}
