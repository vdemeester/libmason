package libmason

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/vdemeester/libmason/test"
)

type WaitClient struct {
	test.NopClient
	success  bool
	duration time.Duration
}

func (c *WaitClient) ContainerWait(ctx context.Context, container string) (int, error) {
	if c.success {
		time.Sleep(c.duration)
		return 0, nil
	}
	return c.NopClient.ContainerWait(ctx, container)
}

func TestContainerWaitErrors(t *testing.T) {
	client := &WaitClient{
		NopClient: test.NopClient{},
	}
	helper := &DefaultHelper{
		client: client,
	}
	_, err := helper.ContainerWait(context.Background(), "container_id", 1*time.Millisecond)
	if err == nil {
		t.Fatalf("expected an error, got nothing")
	}
}

func TestContainerWaitTimesOut(t *testing.T) {
	client := &WaitClient{
		NopClient: test.NopClient{},
		success:   true,
		duration:  100 * time.Millisecond,
	}
	helper := &DefaultHelper{
		client: client,
	}
	_, err := helper.ContainerWait(context.Background(), "container_id", 10*time.Millisecond)
	if err == nil || err.Error() != "Container container_id didn't stop in the specified time : 10ms" {
		t.Fatalf("expected an error, got %v", err)
	}
}

func TestContainerWait(t *testing.T) {
	client := &WaitClient{
		NopClient: test.NopClient{},
		success:   true,
		duration:  2 * time.Millisecond,
	}
	helper := &DefaultHelper{
		client: client,
	}
	statusCode, err := helper.ContainerWait(context.Background(), "container_ud", 10*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if statusCode != 0 {
		t.Fatalf("expected a statusCode 0, got %d", statusCode)
	}
}
