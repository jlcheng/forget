package db

import (
	"fmt"
	"github.com/jlcheng/forget/trace"
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
			trace.Debug(fmt.Sprintf("no-index: .git %v", path))
			return false
		}
	}

	// File name must contain a '.'
	if (!info.IsDir() && strings.LastIndexByte(info.Name(), '.') == -1) {
		trace.Debug(fmt.Sprintf("no-index: name %v", path))
		return false
	}

	// Exclude org mode archive files
	if strings.HasSuffix(info.Name(), ".org_archive") {
		trace.Debug(fmt.Sprintf("no-index: name %v", path))
		return false
	}

	// Omitting large files
	const ONE_MB = int64(1024 * 1024)
	if info.Size() > ONE_MB {
		trace.Debug(fmt.Sprintf("no-index: size %v", path))
		return false
	}
	return true
}
