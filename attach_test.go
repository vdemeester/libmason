package libmason

import (
	"bytes"
	"net/http/httputil"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/vdemeester/libmason/test"
)

type AttachClient struct {
	test.NopClient
	errAttach error
	success   bool
}

func (c *AttachClient) ContainerAttach(ctx context.Context, container string, options types.ContainerAttachOptions) (types.HijackedResponse, error) {
	if c.errAttach != nil || c.success {
		resp := types.HijackedResponse{}

		return resp, c.errAttach
	}
	return c.NopClient.ContainerAttach(ctx, container, options)
}

func TestContainerAttachErrors(t *testing.T) {
	client := &AttachClient{
		NopClient: test.NopClient{},
	}
	helper := &DefaultHelper{
		client: client,
	}
	var b bytes.Buffer
	err := helper.ContainerAttach(context.Background(), "container_id", strings.NewReader(""), &b, &b)
	if err == nil {
		t.Fatalf("expected an error, got nothing")
	}
}

func TestContainerAttachErrNotPersistEOFError(t *testing.T) {
	client := &AttachClient{
		NopClient: test.NopClient{},
		errAttach: httputil.ErrLineTooLong,
	}
	helper := &DefaultHelper{
		client: client,
	}

	var b bytes.Buffer
	err := helper.ContainerAttach(context.Background(), "container_id", strings.NewReader(""), &b, &b)
	if err == nil || err != httputil.ErrLineTooLong {
		t.Fatalf("expected an ErrLineTooLong error, got %v", err)
	}
}

// FIXME(vdemeester) Test more ContainerAttach once holdhijackedconnection is well understood :D
