package db

import (
	"os"
	"path/filepath"
	"strings"
)

// FilterFile returns true for each path that is eligible to be indexed
func FilterFile(path string, info os.FileInfo) bool {
	// Do not index files under .git
	for tmpPath := path;
		tmpPath != "." && tmpPath != string(filepath.Separator);  // invariant
	    tmpPath = filepath.Dir(tmpPath) {
		if strings.HasSuffix(tmpPath, ".git") {
			return false
		}
	}

	// File name must contain a '.'
	if (!info.IsDir() && strings.LastIndexByte(info.Name(), '.') == -1) {
		return false
	}

	// Exclude org mode archive files
	if strings.HasSuffix(info.Name(), ".org_archive") {
		return false
	}

	// Omitting large files
	const ONE_MB = int64(1024 * 1024)
	if info.Size() > ONE_MB {
		return false
	}
	return true
}
