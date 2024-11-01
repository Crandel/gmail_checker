package file

import (
	"fmt"
	"os"

	"github.com/Crandel/gmail/internal/env"
)

const (
	defaultCacheDir = "~/.cache"
	dir             = "mail"
)

func GetCacheDir() (string, error) {
	cacheDir := env.GetEnv("XDG_CACHE_HOME", defaultCacheDir)
	mailCacheDir := fmt.Sprintf("%s/%s", cacheDir, dir)
	err := CreateDirectory(mailCacheDir)
	if err != nil {
		return "", err
	}
	return cacheDir, nil
}

func CreateDirectory(directory string) error {
	if _, err := os.Stat(directory); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(directory, 0755)
}
