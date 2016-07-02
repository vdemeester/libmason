// Package libmason provides the base helper for building client-side builder.
// It consists of an interface (so that you could use the same helper but with a
// different backend) and implementation and utils function.
package libmason

import (
	"io"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
)

var _ Helper = &DefaultHelper{}

// DockerClient defines methods a docker client should provide to libmason.
// It's simply Container & Image methods from docker client.APIClient.
type DockerClient interface {
	client.ContainerAPIClient
	client.ImageAPIClient
}

// DefaultHelper is a client-side builder base helper implementation.
type DefaultHelper struct {
	client       DockerClient
	outputWriter io.Writer
}

// NewHelper creates a new Helper from a docker client
func NewHelper(cli DockerClient) *DefaultHelper {
	return &DefaultHelper{
		client:       cli,
		outputWriter: os.Stdout,
	}
}

// WithOutputWriter lets you specify a writer for the small amount of output this
// package will generate (Pull & such)
func (h *DefaultHelper) WithOutputWriter(w io.Writer) *DefaultHelper {
	h.outputWriter = w
	return h
}

// Helper abstracts calls to a Docker Daemon for a client-side builder.
// It is based on a extraction "revisited" of builder/builder.go from the docker project,
// and define methods a client-side builder might need.
type Helper interface {
	// GetImage looks up a Docker image referenced by `name` and pull it if needed.
	GetImage(ctx context.Context, name string, options types.ImagePullOptions) (types.ImageInspect, error)

	// TagImage tags an image with newTag
	TagImage(ctx context.Context, image string, newReference string) error

	// ContainerCreate creates a new Docker container and returns potential warnings
	ContainerCreate(ctx context.Context, config types.ContainerCreateConfig) (types.ContainerCreateResponse, error)

	// ContainerAttach attaches to container.
	ContainerAttach(ctx context.Context, container string, stdin io.Reader, stdout, stderr io.Writer) error

	// ContainerStart starts a new container
	ContainerStart(ctx context.Context, container string) error

	// ContainerKill stops the container execution abruptly.
	ContainerKill(ctx context.Context, container string) error

	// ContainerRm removes a container specified by `id`.
	ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error

	// Commit creates a new Docker image from an existing Docker container.
	ContainerCommit(ctx context.Context, container string, options types.ContainerCommitOptions) (string, error)

	// ContainerWait stops processing until the given container is stopped.
	ContainerWait(ctx context.Context, container string, timeout time.Duration) (int, error)

	// CopyToContainer copies/extracts a source FileInfo to a destination path inside a container
	// specified by a container object.
	CopyToContainer(ctx context.Context, container string, destPath, srcPath string, decompress bool) error
}
