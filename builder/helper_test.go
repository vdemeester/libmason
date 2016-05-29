package builder

import (
	"errors"
	"io"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/vdemeester/libmason"
)

var (
	errNoHelper = errors.New("No-op helper")
	errStep     = errors.New("error step")

	_ libmason.Helper = &NopHelper{}
)

type NopHelper struct{}

func (h *NopHelper) GetImage(ctx context.Context, name string, options types.ImagePullOptions) (types.ImageInspect, error) {
	return types.ImageInspect{}, errNoHelper
}

func (h *NopHelper) TagImage(ctx context.Context, image string, newReference string) error {
	return errNoHelper
}

func (h *NopHelper) ContainerCreate(ctx context.Context, config types.ContainerCreateConfig) (types.ContainerCreateResponse, error) {
	return types.ContainerCreateResponse{}, errNoHelper
}

func (h *NopHelper) ContainerAttach(ctx context.Context, container string, stdin io.Reader, stdout, stderr io.Writer) error {
	return errNoHelper
}

func (h *NopHelper) ContainerStart(ctx context.Context, containerID string) error {
	return errNoHelper
}

func (h *NopHelper) ContainerKill(ctx context.Context, containerID string) error {
	return errNoHelper
}

func (h *NopHelper) ContainerRemove(ctx context.Context, name string, options types.ContainerRemoveOptions) error {
	return errNoHelper
}

func (h *NopHelper) ContainerCommit(ctx context.Context, name string, options types.ContainerCommitOptions) (string, error) {
	return "", errNoHelper
}

func (h *NopHelper) ContainerWait(ctx context.Context, containerID string, timeout time.Duration) (int, error) {
	return 0, errNoHelper
}

func (h *NopHelper) CopyToContainer(ctx context.Context, container string, destPath, srcPath string, decompress bool) error {
	return errNoHelper
}

type ErrorStep struct{}

func (s *ErrorStep) String() string {
	return "Error"
}

func (s *ErrorStep) Execute(ctx context.Context, helper libmason.Helper, config *Config) (*Config, error) {
	return nil, errStep
}

func stringSliceEqual(s1, s2 []string) bool {
	for i, v := range s1 {
		if s2[i] != v {
			return false
		}
	}
	return true
}
