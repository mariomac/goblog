package assets

import (
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"os"
	"path"
	"strings"
)

const dirStatic = "static"

// We assume fileAssets fit well in memory.
// TODO: for very big assets (e.g. videos) generate a tooBig error and use a normal http.FileServer
type FileAssetGenerator struct {
	rootPath string
}

func (f *FileAssetGenerator) Get(urlPath string) (*WebAsset, error) {
	// replacing "/" by system separator is only needed for Windows:
	// todo: isolate in a windows-only function
	relPath := strings.Split(urlPath[len(pathStatic):], "/")
	absPath := path.Join(append([]string{f.rootPath, dirStatic}, relPath...)...)

	fileBytes, err := os.ReadFile(absPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errNotFound{url: urlPath}
		} else {
			return nil, internalError{cause: fmt.Errorf("reading file %q: %w", urlPath, err)}
		}
	}
	return &WebAsset{
		MimeType: mime.TypeByExtension(path.Ext(urlPath)),
		Body:     fileBytes,
	}, nil
}
