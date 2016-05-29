package builder

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
)

type RemoveHelper struct {
	NopHelper
	success bool
}

func (h *RemoveHelper) ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error {
	if h.success {
		return nil
	}
	return h.NopHelper.ContainerRemove(ctx, container, options)
}

func TestStepRemoveString(t *testing.T) {
	step := WithRemove(&ErrorStep{}).(*RemoveStep)
	actual := step.String()
	expected := "Error (with remove)"
	if actual != expected {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func TestRemoveStepDelegateError(t *testing.T) {
	helper := &RemoveHelper{}
	step := WithRemove(&ErrorStep{})
	_, err := step.Execute(context.Background(), helper, &Config{})
	if err == nil || err != errStep {
		t.Fatalf("expected an errStep error, got %v", err)
	}
}

func TestRemoveStepNoContainerIDError(t *testing.T) {
	helper := &RemoveHelper{}
	step := WithRemove(&NoopStep{})
	_, err := step.Execute(context.Background(), helper, &Config{})
	if err == nil {
		t.Fatalf("expected an error, got nothing")
	}
}

func TestRemoveStepHelperError(t *testing.T) {
	helper := &RemoveHelper{}
	step := WithRemove(&NoopStep{})
	_, err := step.Execute(context.Background(), helper, &Config{
		attributes: map[string]interface{}{
			ContainerID: "ID",
		},
	})
	if err == nil || err != errNoHelper {
		t.Fatalf("expected an errNoHelper, got %v", err)
	}
}

func TestRemoveStep(t *testing.T) {
	helper := &RemoveHelper{
		success: true,
	}
	step := WithRemove(&NoopStep{})
	config, err := step.Execute(context.Background(), helper, &Config{
		attributes: map[string]interface{}{
			ContainerID: "ID",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if v, _ := config.Get(ContainerID); v != "" {
		t.Fatalf("expected ContainerID to be empty, got %s", v)
	}
}
