package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
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

func grep(str string) string {
	r, _ := regexp.Compile(`<fullcount>(.*?)</fullcount>`)
	substring := r.FindString(str)
	re, _ := regexp.Compile(`[\d]`)
	return re.FindString(substring)
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
	body, err := ioutil.ReadAll(resp.Body)
	check(err, "ioutil")
	count := grep(string(body))
	fmt.Printf("%[1]v:%[2]v", configuration.Short, count)
}
