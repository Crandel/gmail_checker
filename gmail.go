package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io"
	//"io/ioutil"
	//"net/http"
	"os"
	//"strings"
	//"regexp"
)

type Label struct {
	Label    string `json:"label,omitempty"`
	Short    string `json:"short_conky,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Configuration struct {
	Email *Label `json:"email"`
}

func check(e error, str string) {
	if e != nil {
		fmt.Printf("%s, %s\n", e, str)
		panic(e)
		os.Exit(1)
	}
}

func readSettings(str string) string {
	filename := fmt.Sprintf("%s/.gmail.json", os.Getenv("HOME"))
	fmt.Println(filename)
	content, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%s, openError\n", err)
		fmt.Println("pfff\n")
		f, err := os.Create(filename)
		check(err, "create")
		defer f.Close()
		ExampleLabel := &Label{
			Label:    "exlabel",
			Short:    "short",
			Username: "Username",
			Password: "Password",
		}
		ExampleConf := &Configuration{
			Email: ExampleLabel,
		}
		example_json, err := json.Marshal(ExampleConf)
		fmt.Println(string(example_json))
		check(err, "Marshal")
		f.WriteString(string(example_json))
		//check(err, "Write")
		return "A file ~/.gmail.json created, please fill username and password"
	} else {
		params := json.NewDecoder(content)
		configuration := Configuration{}
		err := params.Decode(&configuration)
		check(err, "Decode")
		fmt.Printf("%s\n", configuration.Email.Short)
		return str
	}

}

func grep(str io.Reader) {
	doc, err := html.Parse(str)
	check(err, "Parse")
	var f func(*html.Node, bool)
	f = func(n *html.Node, printText bool) {
		if printText && n.Type == html.TextNode {
			fmt.Printf("%q\n", n.Data)
		}
		printText = printText || (n.Type == html.ElementNode && n.Data == "fullcount")
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, printText)
		}
	}
	f(doc, false)
	//re := regexp.MustCompile("<fullcount>(.*?)</fullcount>")
	//fmt.Println(re.FindString(str))
}

func main() {
	base_url := "https://%s:%s@mail.google.com/mail/feed/atom"
	readSettings(base_url)
	//resp, err := http.Get(base_url)
	//check(err, "Get")
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Printf("%s", err)
	//	os.Exit(1)
	//}
	//grep(resp.Body)
	//fmt.Printf(resp.Body)
}
