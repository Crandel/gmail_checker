package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Crandel/gmail/internal/accounts"
)

const fileName = "config.json"
const dir = "mail"

var configDir = fmt.Sprintf("%s/%s", os.Getenv("XDG_CONFIG_HOME"), dir)
var filename = fmt.Sprintf("%s/%s", configDir, fileName)

func GetAccounts() accounts.ListAccounts {
	content, err := os.ReadFile(filename)
	listAccounts := accounts.ListAccounts{}
	if err == nil {
		lAccs := &listAccounts
		err := json.Unmarshal(content, lAccs)
		if err != nil {
			slog.Debug("error during Unmarshal", slog.Any("error", err))
			return listAccounts
		}
	}
	return listAccounts
}

func GetAccount(clientID string) *accounts.Account {
	listAccounts := GetAccounts()
	for _, acc := range listAccounts {
		if acc.ClientID == clientID {
			return &acc
		}
	}
	return nil
}

func AddToConfig() {
	listAccounts := accounts.ListAccounts{}
	if _, err := os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(configDir, os.ModePerm)
		if err != nil {
			slog.Debug("Create config directory", slog.Any("error", err))
		}
	}

	// if file with configuration doesn`t exists this part will create it
	f, err := os.Open(filename)
	if errors.Is(err, os.ErrNotExist) {
		f, err = os.Create(filename)
		if err != nil {
			slog.Debug("error during creation file")
		}
	} else {
		// If file already exists, read it's content
		err = json.NewDecoder(f).Decode(&listAccounts)
		if err != nil {
			slog.Debug("error during Unmarshal", slog.Any("error", err))
		}
	}
	defer f.Close()
	newAccount, err := addNewUser()
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return
	}
	listAccounts = append(listAccounts, newAccount)
	var newJSON []byte
	newJSON, err = json.Marshal(listAccounts)
	if err != nil {
		slog.Debug("error during marshalling", slog.Any("error", err))
	}
	_, err = f.WriteString(string(newJSON))
	if err != nil {
		slog.Debug("error during writing string", slog.Any("error", err))
	}
}

func addNewUser() (accounts.Account, error) {
	// Type necessary account information
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please add mail type. Available types are: gmail, ...")
	mailType, err := reader.ReadString('\n')
	mailType = strings.TrimSuffix(mailType, "\n")
	if err != nil {
		return accounts.Account{}, err
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
	fmt.Println("Please add oauth2 clientSecret")
	clientSecret, err := reader.ReadString('\n')
	clientId = strings.Trim(clientId, "\n")
	if err != nil {
		return accounts.Account{}, err
	}

	return accounts.Account{
		MailType:     mailType,
		Email:        email,
		ClientID:     clientId,
		ClientSecret: clientSecret,
	}, nil
}
