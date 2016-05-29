package builder

import (
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
)

type CommitHelper struct {
	NopHelper
	imageID         string
	expectedChanges []string
}

func (h *CommitHelper) ContainerCommit(ctx context.Context, name string, options types.ContainerCommitOptions) (string, error) {
	if h.imageID != "" {
		if len(h.expectedChanges) != 0 {
			if !reflect.DeepEqual(options.Changes, h.expectedChanges) {
				return "", fmt.Errorf("changes expected %v, got %v", h.expectedChanges, options.Changes)
			}
		}
		return h.imageID, nil
	}
	return h.NopHelper.ContainerCommit(ctx, name, options)
}

func TestStepCommitString(t *testing.T) {
	step := WithCommit(&ErrorStep{}).(*CommitStep)
	actual := step.String()
	expected := "Error (with commit)"
	if actual != expected {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func TestCommitStepDelegateError(t *testing.T) {
	helper := &CommitHelper{}
	step := WithCommit(&ErrorStep{})
	_, err := step.Execute(context.Background(), helper, &Config{})
	if err == nil || err != errStep {
		t.Fatalf("expected an errStep error, got %v", err)
	}
}

func TestCommitStepNoContainerIDError(t *testing.T) {
	helper := &CommitHelper{}
	step := WithCommit(&NoopStep{})
	_, err := step.Execute(context.Background(), helper, &Config{})
	if err == nil {
		t.Fatalf("expected an error, got nothing")
	}
}

func TestCommitStepHelperError(t *testing.T) {
	helper := &CommitHelper{}
	step := WithCommit(&NoopStep{})
	_, err := step.Execute(context.Background(), helper, &Config{
		attributes: map[string]interface{}{
			ContainerID: "ID",
		},
	})
	if err == nil || err != errNoHelper {
		t.Fatalf("expected an errNoHelper, got %v", err)
	}
}

func TestCommitStep(t *testing.T) {
	cmd := []string{"cmd", "arg"}
	entrypoint := []string{"entrypoint", "arg"}
	helper := &CommitHelper{
		imageID: "image_id",
		expectedChanges: []string{
			fmt.Sprintf("CMD %v", cmd),
			fmt.Sprintf("ENTRYPOINT %v", entrypoint),
		},
	}
	step := WithCommit(&NoopStep{})
	config, err := step.Execute(context.Background(), helper, &Config{
		Cmd:        cmd,
		Entrypoint: entrypoint,
		attributes: map[string]interface{}{
			ContainerID: "ID",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if config.ImageID != "image_id" {
		t.Fatalf("expected 'image_id', got %s", config.ImageID)
	}
}
