package utils

import "os"

func CheckFileExists(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil || fileInfo.Name() == "" {
		return false
	}
	return true
}
