package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/Crandel/gmail/internal/accounts"
	"github.com/Crandel/gmail/internal/env"
	"github.com/Crandel/gmail/internal/file"
)

const (
	dir      = "mail"
	fileName = "config.json"
)

var (
	configDir = fmt.Sprintf("%s/%s", env.GetEnv("XDG_CONFIG_HOME", "~/.config"), dir)
	filename  = fmt.Sprintf("%s/%s", configDir, fileName)

	// ErrMailType is an config error type.
	ErrMailType = errors.New("mail type should be only gmail")
)

type inputReader interface {
	ReadString(delim byte) (string, error)
}

// GetFile return io.Reader for config file.
func GetFile() (io.Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// GetAccounts return account list from config file.
func GetAccounts(reader io.Reader) (accounts.ListAccounts, error) {
	listAccounts := accounts.ListAccounts{}
	err := json.NewDecoder(reader).Decode(&listAccounts)
	if err != nil {
		slog.Debug("error during Unmarshal", slog.Any("error", err))
		return listAccounts, err
	}
	return listAccounts, err
}

func readExistingConfig(filename string) (accs accounts.ListAccounts, err error) {
	origFile, err := os.Open(filename)
	if err != nil {
		slog.Error("error during opening file", slog.Any("error", err))
		return nil, err
	}
	defer func() {
		if err = origFile.Close(); err != nil {
			slog.Error("error during closing file", slog.Any("error", err))
		}
	}()

	accs, err = GetAccounts(origFile)
	if err != nil {
		slog.Error("error during Unmarshal", slog.Any("error", err))
		return nil, err
	}

	return accs, nil
}

// AddToConfig adds new user to existed config file.
func AddToConfig() error {
	err := file.CreateDirectory(configDir)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)
	listAccounts, err := readExistingConfig(filename)
	if err != nil {
		return err
	}
	newAccount, err := addNewUser(reader)
	if err != nil {
		slog.Error("An error occurred while reading input. Please try again", slog.Any("error", err))
		return err
	}
	listAccounts = append(listAccounts, newAccount)
	slog.Debug("Accounts list", slog.Any("list", listAccounts))

	var newJSON []byte
	newJSON, err = json.MarshalIndent(listAccounts, "", "  ")
	if err != nil {
		slog.Error("error during marshaling", slog.Any("error", err))
		return err
	}
	err = os.WriteFile(filename, newJSON, 0666)
	if err != nil {
		slog.Error("error during writing string", slog.Any("error", err))
		return err
	}

	return err
}

func addNewUser(reader inputReader) (accounts.Account, error) {
	fmt.Println("Please add short alias for this account")
	mailAlias, err := reader.ReadString('\n')
	mailAlias = strings.TrimSpace(mailAlias)
	if err != nil {
		return accounts.Account{}, err
	}

	fmt.Println("Please add mail type. Available types are: gmail")
	mailType, err := reader.ReadString('\n')
	mailType = strings.TrimSpace(mailType)
	if err != nil {
		return accounts.Account{}, err
	}
	mailT := accounts.MailType(mailType)
	if mailT != accounts.Gmail {
		return accounts.Account{}, ErrMailType
	}

	fmt.Println("Please add email address")
	email, err := reader.ReadString('\n')
	email = strings.TrimSpace(email)
	if err != nil {
		return accounts.Account{}, err
	}

	fmt.Println("Please add oauth2 clientId")
	clientId, err := reader.ReadString('\n')
	clientId = strings.TrimSpace(clientId)
	if err != nil {
		return accounts.Account{}, err
	}

	return accounts.Account{
		Short:    mailAlias,
		MailType: mailT,
		Email:    email,
		ClientID: clientId,
	}, nil
}
