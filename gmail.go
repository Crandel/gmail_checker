package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"text/template"
)

type Account struct {
	Account  string `json:"account"`
	Short    string `json:"short_conky"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func check(e error, str string) {
	if e != nil {
		fmt.Printf("%s, %s\n", e, str)
		panic(e)
		os.Exit(1)
	}
}

func readSettings() Account {
	filename := fmt.Sprintf("%s/.gmail.json", os.Getenv("HOME"))
	content, err := os.Open(filename)
	if err != nil {
		f, err := os.Create(filename)
		check(err, "create")
		defer f.Close()
		exampleAccount := Account{
			Account:  "ACCOUNT",
			Short:    "SHORT",
			Username: "USERNAME",
			Password: "PASSWORD",
		}
		example_json, err := json.Marshal(exampleAccount)
		check(err, "Marshal")
		f.WriteString(string(example_json))
		check(err, "Write")
		return exampleAccount
	} else {
		params := json.NewDecoder(content)
		configuration := Account{}
		err := params.Decode(&configuration)
		check(err, "Decode")

		return configuration
	}
}

func grep(str io.Reader) {
	doc, err := html.Parse(str)
	check(err, "Parse")
	var f func(*html.Node, bool)
	f = func(n *html.Node, printText bool) {
		if printText && n.Type == html.TextNode {
			fmt.Println(n.Data)
		}
		printText = printText || (n.Type == html.ElementNode && n.Data == "fullcount")
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, printText)
		}
	}
	f(doc, false)
}

func main() {
	base_url := "https://{{.Username}}:{{.Password}}@mail.google.com/mail/feed/atom"
	configuration := readSettings()
	t := template.New(configuration.Account)
	t, _ = t.Parse(base_url)
	buf := new(bytes.Buffer)
	t.Execute(buf, configuration)
	resp, err := http.Get(buf.String())
	check(err, "Get")
	grep(resp.Body)
}
