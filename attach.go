package libmason

import (
	"io"
	"net/http/httputil"

	"golang.org/x/net/context"

	// FIXME(vdemeester) Remove dependency for docker/docker
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/types"
)

// ContainerAttach attaches to container.
func (h *DefaultHelper) ContainerAttach(ctx context.Context, container string, stdin io.Reader, stdout, stderr io.Writer) error {
	// pipe stdin, stderr and stdout (and stream) in containerAttachOptions
	resp, errAttach := h.client.ContainerAttach(ctx, container, types.ContainerAttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: true, // Should it be an option ?
	})
	if errAttach != nil && errAttach != httputil.ErrPersistEOF {
		return errAttach
	}
	defer resp.Close()

	tty := true // FIXME(vdemeester) TTY stuff
	if err := holdHijackedConnection(ctx, tty, stdin, stdout, stderr, resp); err != nil {
		return err
	}

	if errAttach != nil {
		return errAttach
	}
	return nil
}

// FIXME(vdemeester) Handle context :)
func holdHijackedConnection(ctx context.Context, tty bool, inputStream io.Reader, outputStream, errorStream io.Writer, resp types.HijackedResponse) error {
	var err error

	receiveStdout := make(chan error, 1)
	if outputStream != nil || errorStream != nil {
		go func() {
			// When TTY is ON, use regular copy
			if tty && outputStream != nil {
				_, err = io.Copy(outputStream, resp.Reader)
			} else {
				_, err = stdcopy.StdCopy(outputStream, errorStream, resp.Reader)
			}

			receiveStdout <- err
		}()
	}

	stdinDone := make(chan struct{})
	go func() {
		if inputStream != nil {
			io.Copy(resp.Conn, inputStream)
		}

		if err := resp.CloseWrite(); err != nil {
			// FIXME(vdemeester) log something ?
		}
		close(stdinDone)
	}()

	select {
	case err := <-receiveStdout:
		if err != nil {
			return err
		}
	case <-stdinDone:
		if outputStream != nil || errorStream != nil {
			select {
			case err := <-receiveStdout:
				if err != nil {
					return err
				}
			case <-ctx.Done():
			}
		}
	case <-ctx.Done():
	}

	return nil
}
