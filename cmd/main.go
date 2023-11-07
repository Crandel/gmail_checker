package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/Crandel/gmail/internal/accounts"
	"github.com/Crandel/gmail/internal/logging"
	"github.com/Crandel/gmail/internal/mails"
)

func getAccounts() accounts.ListAccounts {
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

func main() {
	debug := os.Getenv("DEBUG")
	logLevel := slog.LevelInfo
	showSources := false
	if debug == "1" {
		showSources = true
		logLevel = slog.LevelDebug
	}

	logging.InitLogger(logLevel, showSources)
	// Check if domain online
	resp, err := http.Get("https://mail.google.com")
	if err == nil || resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMovedPermanently {
		channel := make(chan string)
		defer close(channel)
		listAccounts := getAccounts()
		for _, acc := range listAccounts {
			// separate all network requests to goroutines
			go mails.GetMailCount(channel, acc)
		}
		accLen := len(listAccounts)
		counts := make([]string, accLen)
		for i := 0; i < accLen; i++ {
			counts[i] = <-channel
		}
		fmt.Println(strings.Join(counts, ""))
	}
}
