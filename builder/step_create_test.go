package builder

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
)

type CreateHelper struct {
	NopHelper
	containerID        string
	expectedCmd        []string
	expectedEntrypoint []string
	expectedStdin      bool
}

func (h *CreateHelper) ContainerCreate(ctx context.Context, config types.ContainerCreateConfig) (types.ContainerCreateResponse, error) {
	if h.containerID != "" {
		if len(h.expectedCmd) != 0 {
			if !stringSliceEqual(config.Config.Cmd, h.expectedCmd) {
				return types.ContainerCreateResponse{}, fmt.Errorf("cmd expected %v, got %v", h.expectedCmd, config.Config.Cmd)
			}
		}
		if len(h.expectedEntrypoint) != 0 {
			if !stringSliceEqual(config.Config.Entrypoint, h.expectedEntrypoint) {
				return types.ContainerCreateResponse{}, fmt.Errorf("entrypoint expected %v, got %v", h.expectedEntrypoint, config.Config.Entrypoint)
			}
		}
		if h.expectedStdin != config.Config.StdinOnce {
			return types.ContainerCreateResponse{}, fmt.Errorf("stdin expected %v, got %v", h.expectedStdin, config.Config.StdinOnce)
		}
		return types.ContainerCreateResponse{
			ID: h.containerID,
		}, nil
	}
	return h.NopHelper.ContainerCreate(ctx, config)
}

func TestStepCreateString(t *testing.T) {
	step := WithCreate(&ErrorStep{}, []string{}, []string{}, false).(*CreateStep)
	actual := step.String()
	expected := "Error (with create: [] [])"
	if actual != expected {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func TestCreateStepHelperError(t *testing.T) {
	helper := &CreateHelper{}
	step := WithCreate(&NoopStep{}, []string{}, []string{}, false)
	_, err := step.Execute(context.Background(), helper, &Config{})
	if err == nil || err != errNoHelper {
		t.Fatalf("expected an errNoHelper, got %v", err)
	}
}

func TestCreateStepDelegateError(t *testing.T) {
	helper := &CreateHelper{
		containerID: "container_id",
	}
	step := WithCreate(&ErrorStep{}, []string{}, []string{}, false)
	_, err := step.Execute(context.Background(), helper, &Config{})
	if err == nil || err != errStep {
		t.Fatalf("expected an errStep error, got %v", err)
	}
}

func TestCreateStep(t *testing.T) {
	cmd := []string{"cmd", "arg"}
	entrypoint := []string{"entrypoint", "arg"}
	helper := &CreateHelper{
		containerID:        "container_id",
		expectedCmd:        cmd,
		expectedEntrypoint: entrypoint,
		expectedStdin:      true,
	}
	step := WithCreate(&NoopStep{}, entrypoint, cmd, true)
	config, err := step.Execute(context.Background(), helper, &Config{})
	if err != nil {
		t.Fatal(err)
	}
	actual, _ := config.Get(ContainerID)
	if actual != "container_id" {
		t.Fatalf("expected 'container_id', got %v", actual)
	}
}
