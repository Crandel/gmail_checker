package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// ListAccounts - list of accounts from config file
var ListAccounts []Account

// Account type - description of account
type Account struct {
	MailType string `json:"mail_type"`
	Account  string `json:"account"`
	Short    string `json:"short_conky"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Feed struct {
	XMLName   xml.Name `xml."feed"`
	Xmlns     string   `xml."xmlns,attr"`
	Version   string   `xml."version,attr"`
	Title     string   `xml."name"`
	Tagline   string   `xml."tagline"`
	Fullcount string   `xml."fullcount"`
	Link      Link     `xml."link"`
}

type Link struct {
	XMLName xml.Name `xml."link"`
	Rel     string   `xml."rel,attr"`
	Href    string   `xml."href,attr"`
	Type    string   `xml."type,attr"`
}

// Error function
func check(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}

func init() {
	filename := fmt.Sprintf("%s/.email.json", os.Getenv("HOME"))
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		// if file with configuration does`nt exists this part will create it
		f, err := os.Create(filename)
		check(err)
		defer f.Close()
		exampleAccount := Account{
			MailType: "gmail",
			Account:  "ACCOUNT",
			Short:    "SHORT",
			Email:    "EMAIL@gmail.com",
			Password: "PASSWORD",
		}
		listAccounts := []Account{exampleAccount}
		exampleJSON, err := json.Marshal(listAccounts)
		check(err)
		f.WriteString(string(exampleJSON))
		ListAccounts = listAccounts
	} else {
		listAccounts := &ListAccounts
		err := json.Unmarshal(content, listAccounts)
		check(err)
	}
}

func extract_count(str string) string {
	var feed Feed
	xml.Unmarshal([]byte(str), &feed)
	return feed.Fullcount
}

// getMailCount - new goroutine for checking emails
func getMailCount(channel chan<- string, acc Account) {
	url := "https://mail.google.com/mail/feed/atom"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	check(err)
	tokenStr := fmt.Sprintf("%s:%s", acc.Email, acc.Password)
	token := base64.StdEncoding.EncodeToString([]byte(tokenStr))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
	resp, err := client.Do(req)
	check(err)
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	count := extract_count(string(body))
	channel <- fmt.Sprintf("%[1]v:%[2]v ", acc.Short, count)
}

func main() {
	// Check if domain online
	resp, err := http.Get("https://mail.google.com")
	if err == nil || resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMovedPermanently {
		channel := make(chan string)
		defer close(channel)
		for _, acc := range ListAccounts {
			// separate all network requests to goroutines
			go getMailCount(channel, acc)
		}
		accLen := len(ListAccounts)
		counts := make([]string, accLen)
		for i := 0; i < accLen; i++ {
			counts[i] = <-channel
		}
		fmt.Println(strings.Join(counts, ""))
	}
}
