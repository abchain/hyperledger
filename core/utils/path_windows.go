package utils

import (
	"golang.org/x/sys/windows/registry"
	"path/filepath"
	"strings"
)

func CanonicalizePath(path string) string {

	path = filepath.ToSlash(path)
	pathexpanded, err := registry.ExpandString(path)
	if err == nil {
		path = pathexpanded
	}

	if !strings.HasSuffix(path, "\\") {
		path = path + "\\"
	}

	return path

}

func CanonicalizeFilePath(path string) string {

	path = filepath.ToSlash(path)
	pathpathexpanded, err := registry.ExpandString(path)
	if err == nil {
		return pathpathexpanded
	} else {
		return path
	}

}
