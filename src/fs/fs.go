// File System tools
package fs

import (
	"regexp"
	"path/filepath"
	"os"
)

// Returns the paths of all the files contained in the folder and subfolders of the path whose
// file name matches with the regular expression.
// The paths are returned in alphabetical order.
// It excludes the directories.
func Search(folder string, regexp *regexp.Regexp) []string {
	paths := make([]string, 0, 32)
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		// If there is any error, just ignores the file
		if err == nil && !info.IsDir() && regexp.MatchString(info.Name()) {
			paths = append(paths, path)
		}
		return nil
	})
	return paths
}
