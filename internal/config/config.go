package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/Crandel/gmail/internal/accounts"
)

func GetAccounts() accounts.ListAccounts {
	filename := fmt.Sprintf("%s/.email.json", os.Getenv("HOME"))
	content, err := os.ReadFile(filename)
	listAccounts := accounts.ListAccounts{}
	if err != nil {
		// if file with configuration does`nt exists this part will create it
		f, err := os.Create(filename)
		if err != nil {
			slog.Debug("error during creation file")
			return listAccounts
		}
		defer f.Close()
		exampleAccount := accounts.Account{
			MailType: "gmail",
			Account:  "ACCOUNT",
			Short:    "SHORT",
			Email:    "EMAIL@gmail.com",
			Password: "PASSWORD",
		}
		listAccounts = append(listAccounts, exampleAccount)
		var exampleJSON []byte
		exampleJSON, err = json.Marshal(listAccounts)
		if err != nil {
			slog.Debug("error during marshalling", slog.Any("error", err))
		}
		_, err = f.WriteString(string(exampleJSON))
		if err != nil {
			slog.Debug("error during writing string", slog.Any("error", err))
		}
	} else {
		lAccs := &listAccounts
		err := json.Unmarshal(content, lAccs)
		if err != nil {
			slog.Debug("error during Unmarshal", slog.Any("error", err))
			return listAccounts
		}
	}
	return listAccounts
}
