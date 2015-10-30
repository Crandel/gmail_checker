package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	//"net/http"
	"os"
	"strings"
	//	"regexp"
)

func readSettings(str string) string {

	content, err := ioutil.ReadFile(fmt.Sprintf("%s/.gmailrc", os.Getenv("HOME")))
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	params := strings.Split(string(content), ":")
	for i := 0; i < len(params); i++ {
		fmt.Printf("%s", params[i])
	}
	return str
}

func grep(str io.Reader) {
	doc, err := html.Parse(str)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
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
	//if err != nil {
	//	fmt.Printf("%s", err)
	//	os.Exit(1)
	//}
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Printf("%s", err)
	//	os.Exit(1)
	//}
	//grep(resp.Body)
	//fmt.Printf(resp.Body)
}
