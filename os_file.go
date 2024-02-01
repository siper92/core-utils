package core_utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const NoFileHash = "noFile"

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func IsDir(path string) bool {
	if stat, err := os.Stat(path); err != nil {
		ErrorWarning(err)
		return false
	} else {
		return stat.IsDir()
	}
}

func CreateFilepathDir(filePath string) error {
	dirPath := filepath.Dir(filePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetFileContent(path string) (*bytes.Buffer, error) {
	var dat []byte
	var err error

	if FileExists(path) {
		dat, err = os.ReadFile(path)
		ErrorWarning(err)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("file %s is missing", path)
	}

	return bytes.NewBuffer(dat), err
}

func MustGetFileContent(path string) *bytes.Buffer {
	data, err := GetFileContent(path)
	ErrorWarning(err)

	return data
}

func GetFileHash(path string) string {
	if FileExists(path) {
		dat, err := os.ReadFile(path)
		ErrorWarning(err)

		if err == nil {
			return RawSha256(dat)
		}
	}

	Debug("File %s is missing", path)
	return NoFileHash
}

func GetRelativePath(dir string, path string) string {
	var err error

	path, err = filepath.Abs(path)
	StopOnError(err)

	dir, err = filepath.Abs(dir)
	StopOnError(err)

	return strings.TrimLeft(strings.TrimPrefix(path, dir), "/")
}

func WriteToFile(path string, data *bytes.Buffer) error {
	var err error

	// Create file if not exists or truncate
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	// Write data to file
	writtenBytes, err := file.Write(data.Bytes())
	_ = writtenBytes

	return err
}

func AbsPath(path string) string {
	absPath, err := filepath.Abs(path)
	StopOnError(err)

	return absPath
}
