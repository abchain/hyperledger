// +build !windows

package utils

import (
	"path/filepath"
	"strings"
)

func CanonicalizePath(path string) string {

	path = filepath.ToSlash(path)
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	return path

}

func CanonicalizeFilePath(path string) string {

	return filepath.ToSlash(path)

}
