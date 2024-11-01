package keyring

import (
	"log/slog"

	"github.com/zalando/go-keyring"
)

type KeyringHandler interface {
	GetEntry(key string) (string, error)
	SetEntry(key string, data string) error
}
type keyringHandler struct {
	name string
}

func NewKeyring(name string) KeyringHandler {
	return keyringHandler{
		name: name,
	}
}

func (kh keyringHandler) GetEntry(key string) (string, error) {
	// get password
	secret, err := keyring.Get(kh.name, key)
	if err != nil {
		slog.Debug("can't get credentials from keyring", slog.Any("error", err))
		return "", err
	}
	return secret, nil
}

func (kh keyringHandler) SetEntry(key string, data string) error {
	err := keyring.Set(kh.name, key, data)
	if err != nil {
		slog.Debug("can't save credentials to keyring", slog.Any("error", err))
		return err
	}
	return nil
}
