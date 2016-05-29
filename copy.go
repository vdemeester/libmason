package libmason

import (
	"golang.org/x/net/context"

	"github.com/docker/docker/pkg/archive"
	"github.com/docker/engine-api/types"
)

// CopyToContainer copies/extracts a source FileInfo to a destination path inside a container
// specified by a container object.
func (h *DefaultHelper) CopyToContainer(ctx context.Context, container string, destPath, srcPath string, decompress bool) error {
	dstInfo := archive.CopyInfo{Path: destPath}
	// FIXME(vdemeester) handle link follow here ?
	dstStat, err := h.client.ContainerStatPath(ctx, container, destPath)
	if err == nil {
		dstInfo.Exists, dstInfo.IsDir = true, dstStat.Mode.IsDir()
	}
	srcInfo, err := archive.CopyInfoSourcePath(srcPath, false)
	if err != nil {
		return err
	}
	srcArchive, err := archive.TarResource(srcInfo)
	if err != nil {
		return err
	}
	defer srcArchive.Close()
	destDir, preparedArchive, err := archive.PrepareArchiveCopy(srcArchive, srcInfo, dstInfo)
	if err != nil {
		return err
	}
	defer preparedArchive.Close()
	return h.client.CopyToContainer(ctx, container, destDir, preparedArchive, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
}
