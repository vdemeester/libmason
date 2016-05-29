package libmason

import (
	"encoding/base64"
	"encoding/json"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"

	// FIXME(vdemeester) Remove dependency for docker/docker
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
)

// GetImage looks up a Docker image referenced by `name` and pull it if needed.
func (h *DefaultHelper) GetImage(ctx context.Context, ref string, options types.ImagePullOptions) (types.ImageInspect, error) {
	if imageInspect, _, err := h.client.ImageInspectWithRaw(ctx, ref, false); err == nil {
		return imageInspect, nil
	}
	// FIXME(vdemeestr) Handle RegistryAuth (authConfig)
	responseBody, err := h.client.ImagePull(ctx, ref, options)
	if err != nil {
		return types.ImageInspect{}, err
	}

	outFd, isTerminalOut := term.GetFdInfo(h.outputWriter)

	defer responseBody.Close()
	if err := jsonmessage.DisplayJSONMessagesStream(responseBody, h.outputWriter, outFd, isTerminalOut, nil); err != nil {
		return types.ImageInspect{}, err
	}
	imageInspect, _, err := h.client.ImageInspectWithRaw(context.Background(), ref, false)
	return imageInspect, err
}

// encodeAuthToBase64 serializes the auth configuration as JSON base64 payload
func encodeAuthToBase64(authConfig types.AuthConfig) (string, error) {
	buf, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}
