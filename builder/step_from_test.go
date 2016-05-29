package builder

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
)

type FromHelper struct {
	NopHelper
	imageID     string
	expectedRef string
}

func (h *FromHelper) GetImage(ctx context.Context, ref string, options types.ImagePullOptions) (types.ImageInspect, error) {
	if h.imageID != "" {
		if h.expectedRef != ref {
			return types.ImageInspect{}, fmt.Errorf("ref expected %s, got %s", h.expectedRef, ref)
		}
		return types.ImageInspect{
			ID: h.imageID,
		}, nil
	}
	return h.NopHelper.GetImage(ctx, ref, options)
}

func TestFromStepString(t *testing.T) {
	step := &FromStep{
		Reference: "reference",
	}
	actual := step.String()
	expected := "FROM reference"
	if actual != expected {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func TestFromStepWithImageIDNotEmptyError(t *testing.T) {
	helper := &FromHelper{}
	step := &FromStep{
		Reference: "reference",
	}
	_, err := step.Execute(context.Background(), helper, &Config{
		ImageID: "image_id",
	})
	if err == nil {
		t.Fatalf("expected an error, got nothing")
	}
}

func TestFromStepHelperError(t *testing.T) {
	helper := &FromHelper{}
	step := &FromStep{
		Reference: "reference",
	}
	_, err := step.Execute(context.Background(), helper, &Config{})
	if err == nil || err != errNoHelper {
		t.Fatalf("expected an errNoHelper, got %v", err)
	}
}

func TestFromStep(t *testing.T) {
	helper := &FromHelper{
		imageID:     "image_id",
		expectedRef: "reference",
	}
	step := &FromStep{
		Reference: "reference",
	}
	config, err := step.Execute(context.Background(), helper, &Config{})
	if err != nil {
		t.Fatal(err)
	}
	if config.ImageID != "image_id" {
		t.Fatalf("expected 'image_id', got %s", config.ImageID)
	}
}
