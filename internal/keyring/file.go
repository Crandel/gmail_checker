package keyring

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// NewFileKeyring return fileKeyring with corresponding file.
func NewFileKeyring(dir, name string) (Handler, error) {
	filename := fmt.Sprintf("%s/%s_%s", dir, name, fileName)
	f, err := os.Open(filename)
	if err != nil {
		f, err = os.Create(filename)
		if err != nil {
			return nil, err
		}
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}

	return fileKeyring{
		file: filename,
	}, nil
}

func (fk fileKeyring) getMap() (content, error) {
	c := make(map[string]string)
	origFile, err := os.Open(fk.file)
	if err != nil {
		slog.Error("error during opening file", slog.Any("error", err))
		return nil, err
	}
	defer func() {
		if err = origFile.Close(); err != nil {
			slog.Error("error during closing file", slog.Any("error", err))
		}
	}()

	err = json.NewDecoder(origFile).Decode(&c)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			slog.Error("error during Unmarshal", slog.Any("error", err))
			return nil, err
		}
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
	con[key] = data
	slog.Debug("con", slog.Any("map", con))
	newJSON, err = json.MarshalIndent(con, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(fk.file, newJSON, 0666)
	if err != nil {
		slog.Error("error during writing string", slog.Any("error", err))
		return err
	}

	return err
}
