package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"regexp"
	"text/template"
)

// Account type - description of account
type Account struct {
	Account  string `json:"account"`
	Short    string `json:"short_conky"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Error function
func check(e error) {
	if e != nil {
		l3, err := syslog.New(syslog.LOG_ERR, "Go gmail")
		defer l3.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Fatal(e.Error())
	}
}

func readSettings() []Account {
	filename := fmt.Sprintf("%s/.gmail.json", os.Getenv("HOME"))
	content, err := ioutil.ReadFile(filename)
	var listAccounts []Account
	if err != nil {
		// if file with configuration does`nt exists this part will create it
		f, err := os.Create(filename)
		check(err)
		defer f.Close()
		exampleAccount := Account{
			Account:  "ACCOUNT",
			Short:    "SHORT",
			Email:    "EMAIL@gmail.com",
			Password: "PASSWORD",
		}
		exAccounts := []Account{exampleAccount}
		exampleJSON, err := json.Marshal(exAccounts)
		check(err)
		f.WriteString(string(exampleJSON))
		listAccounts = exAccounts
	} else {
		var configuration []Account
		err := json.Unmarshal(content, &configuration)
		check(err)
		listAccounts = configuration
	}
	return listAccounts
}

func grep(str string) string {
	r, _ := regexp.Compile(`<fullcount>(.*?)</fullcount>`)
	substring := r.FindString(str)
	re, _ := regexp.Compile(`[\d+]`)
	return re.FindString(substring)
}

func main() {
	baseURL := "https://{{.Email}}:{{.Password}}@mail.google.com/mail/feed/atom"
	configuration := readSettings()
	for index := range configuration {
		t := template.New(configuration[index].Account)
		t, _ = t.Parse(baseURL)
		buf := new(bytes.Buffer)
		err := t.Execute(buf, configuration[index])
		check(err)
		resp, err := http.Get(buf.String())
		check(err)
		body, err := ioutil.ReadAll(resp.Body)
		check(err)
		count := grep(string(body))
		fmt.Printf("%[1]v:%[2]v ", configuration[index].Short, count)
	}
}
