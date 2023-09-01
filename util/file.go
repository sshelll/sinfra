package util

import (
	"os"
	"time"
)

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDirExist(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}

func LastModTime(path string) *time.Time {
	stat, err := os.Stat(path)
	if err != nil {
		return nil
	}
	t := stat.ModTime()
	return &t
}
