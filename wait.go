package libmason

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
)

type waitReturn struct {
	statusCode int
	err        error
}

// ContainerWait stops processing until the given container is stopped.
func (h *DefaultHelper) ContainerWait(ctx context.Context, container string, timeout time.Duration) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ret := make(chan waitReturn)
	go func() {
		i, err := h.client.ContainerWait(ctx, container)
		ret <- waitReturn{i, err}
	}()

	select {
	case r := <-ret:
		return r.statusCode, r.err
	case <-ctx.Done():
		return -1, fmt.Errorf("Container %s didn't stop in the specified time : %v", container, timeout)
	}
}
