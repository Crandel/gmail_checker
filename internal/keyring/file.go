package keyring

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

const (
	fileName = "keyring.json"
)

type content map[string]string

type fileKeyring struct {
	file string
}

func NewFileKeyring(dir, name string) KeyringHandler {
	filename := fmt.Sprintf("%s/%s_%s", dir, fileName, name)
	return fileKeyring{
		file: filename,
	}
}

func (fk fileKeyring) getMap() (content, error) {
	var c content
	origFile, err := os.Open(fk.file)
	if err != nil {
		slog.Error("error during opening file", slog.Any("error", err))
		return nil, err
	}
	defer origFile.Close()
	err = json.NewDecoder(origFile).Decode(&c)
	if err != nil {
		slog.Error("error during Unmarshal", slog.Any("error", err))
		return nil, err
	}
	return c, nil
}

func (fk fileKeyring) GetEntry(key string) (string, error) {
	con, err := fk.getMap()
	if err != nil {
		return "", err
	}
	entry, ok := con[key]
	// get password
	if !ok {
		return "", err
	}
	return entry, nil
}

func (fk fileKeyring) SetEntry(key string, data string) error {
	var newJSON []byte
	con, err := fk.getMap()
	if err != nil {
		return err
	}
	_, ok := con[key]
	if ok {
		return fmt.Errorf("entry for key: %s alredy exists", key)
	}
	newJSON, err = json.MarshalIndent(con, "", "  ")

	err = os.WriteFile(fk.file, newJSON, 0666)
	if err != nil {
		slog.Error("error during writing string", slog.Any("error", err))
		return err
	}

	return err
}
