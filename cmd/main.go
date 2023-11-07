package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Crandel/gmail/internal/accounts"
)

func init() {
	filename := fmt.Sprintf("%s/.email.json", os.Getenv("HOME"))
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		// if file with configuration does`nt exists this part will create it
		f, err := os.Create(filename)
		check(err)
		defer f.Close()
		exampleAccount := accounts.Account{
			MailType: "gmail",
			Account:  "ACCOUNT",
			Short:    "SHORT",
			Email:    "EMAIL@gmail.com",
			Password: "PASSWORD",
		}
		listAccounts := []accounts.Account{exampleAccount}
		exampleJSON, err := json.Marshal(listAccounts)
		f.WriteString(string(exampleJSON))
		ListAccounts = listAccounts
	} else {
		listAccounts := &ListAccounts
		err := json.Unmarshal(content, listAccounts)
	}
}

func main() {
	// Check if domain online
	resp, err := http.Get("https://mail.google.com")
	if err == nil || resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMovedPermanently {
		channel := make(chan string)
		defer close(channel)
		for _, acc := range ListAccounts {
			// separate all network requests to goroutines
			go extractors.GetMailCount(channel, acc)
		}
		accLen := len(ListAccounts)
		counts := make([]string, accLen)
		for i := 0; i < accLen; i++ {
			counts[i] = <-channel
		}
		fmt.Println(strings.Join(counts, ""))
	}
}
