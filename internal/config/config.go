package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/Crandel/gmail/internal/accounts"
)

const fileName = ".email.json"

var filename = fmt.Sprintf("%s/%s", os.Getenv("HOME"), fileName)

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

func CreateConfig() {
	// if file with configuration doesn`t exists this part will create it
	_, err := os.ReadFile(filename)
	if err == nil {
		// If file already exists, ignore it
		return
	}
	f, err := os.Create(filename)
	if err != nil {
		slog.Debug("error during creation file")
	}
	defer f.Close()
	exampleAccount := accounts.Account{
		MailType: "gmail",
		Account:  "ACCOUNT",
		Short:    "SHORT",
		Email:    "EMAIL@gmail.com",
		Password: "PASSWORD",
	}

	listAccounts := accounts.ListAccounts{exampleAccount}
	var exampleJSON []byte
	exampleJSON, err = json.Marshal(listAccounts)
	if err != nil {
		slog.Debug("error during marshalling", slog.Any("error", err))
	}
	_, err = f.WriteString(string(exampleJSON))
	if err != nil {
		slog.Debug("error during writing string", slog.Any("error", err))
	}
}
