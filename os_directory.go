package core_utils

import (
	"os"
	"path/filepath"
	"strings"
)

func DirectoryExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func GetFilesList(rootPath string) []string {
	return GetFileListExcluding(rootPath, []string{})
}

func GetFileListExcluding(rootPath string, excluded []string) []string {
	if !DirectoryExists(rootPath) {
		Notice("Directory " + rootPath + " does not exist")
		return []string{}
	}

	var filePaths []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			for _, exclude := range excluded {
				if strings.Contains(path, exclude) {
					return nil
				}
			}

			filePaths = append(filePaths, path)
		}

		return nil
	})
	ErrorWarning(err)

	return filePaths
}
