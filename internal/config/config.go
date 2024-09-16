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
)

const fileName = "config.json"
const dir = "mail"

var ConfigDir = fmt.Sprintf("%s/%s", os.Getenv("XDG_CONFIG_HOME"), dir)
var filename = fmt.Sprintf("%s/%s", ConfigDir, fileName)

type inputReader interface {
	ReadString(delim byte) (string, error)
}

func GetFilename() string {
	return filename
}

func GetAccounts(reader io.Reader) (accounts.ListAccounts, error) {
	listAccounts := accounts.ListAccounts{}
	err := json.NewDecoder(reader).Decode(&listAccounts)
	if err == nil {
		slog.Debug("error during Unmarshal", slog.Any("error", err))
		return listAccounts, err
	}
	return listAccounts, err
}

func AddToConfig() error {
	listAccounts := accounts.ListAccounts{}
	if _, err := os.Stat(ConfigDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(ConfigDir, os.ModePerm)
		if err != nil {
			slog.Error("Create config directory", slog.Any("error", err))
			return err
		}
	}

	// if file with configuration doesn`t exists this part will create it
	origFile, err := os.Open(filename)
	if errors.Is(err, os.ErrNotExist) {
		origFile, err = os.Create(filename)
		if err != nil {
			slog.Error("error during creation file", slog.String("filename", filename))
			return err
		}
	}
	// If file already exists, read it's content
	listAccounts, err = GetAccounts(origFile)
	if err != nil {
		slog.Error("error during Unmarshal", slog.Any("error", err))
		return err
	}

	origFile.Close()
	reader := bufio.NewReader(os.Stdin)

	newAccount, err := addNewUser(reader)
	if err != nil {
		slog.Error("An error occured while reading input. Please try again", slog.Any("error", err))
		return err
	}
	listAccounts = append(listAccounts, newAccount)
	slog.Debug("Accounts list", slog.Any("list", listAccounts))
	var newJSON []byte
	newJSON, err = json.Marshal(listAccounts)
	if err != nil {
		slog.Error("error during marshalling", slog.Any("error", err))
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
	// Type necessary account information

	fmt.Println("Please add short alias for this account")
	mailAlias, err := reader.ReadString('\n')
	mailAlias = strings.TrimSuffix(mailAlias, "\n")
	if err != nil {
		return accounts.Account{}, err
	}

	fmt.Println("Please add mail type. Available types are: gmail, ...")
	mailType, err := reader.ReadString('\n')
	mailType = strings.TrimSuffix(mailType, "\n")
	if err != nil {
		return accounts.Account{}, err
	}
	mailT := accounts.MailType(mailType)
	if mailT != accounts.Gmail {
		return accounts.Account{}, errors.New("mail type should be only gmail")
	}

	fmt.Println("Please add email address")
	email, err := reader.ReadString('\n')
	email = strings.Trim(email, "\n")
	if err != nil {
		return accounts.Account{}, err
	}

	fmt.Println("Please add oauth2 clientId")
	clientId, err := reader.ReadString('\n')
	clientId = strings.Trim(clientId, "\n")
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
