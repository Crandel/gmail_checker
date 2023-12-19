package keyring

import (
	"log/slog"

	"github.com/zalando/go-keyring"
)

const service = "gmail_checker"

func GetEntry(key string) (string, error) {
	// get password
	secret, err := keyring.Get(service, key)
	if err != nil {
		slog.Debug("can't get credentials from keyring", slog.Any("error", err))
		return "", err
	}
	return secret, nil
}

func SetEntry(key string, data string) error {
	err := keyring.Set(service, key, data)
	if err != nil {
		slog.Debug("can't save credentials to keyring", slog.Any("error", err))
		return err
	}
	return nil
}
