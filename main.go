package main

import (
	"flag"
	"fmt"
	"github.com/golang-commonmark/markdown"
	"io/ioutil"
	"log"
	"net/http"
)

type snippet struct {
	content string
	lang    string
}

func getSnippet(tok markdown.Token) snippet {
	switch tok := tok.(type) {
	case *markdown.CodeBlock:
		return snippet{
			tok.Content,
			"code",
		}
	case *markdown.CodeInline:
		return snippet{
			tok.Content,
			"code inline",
		}
	case *markdown.Fence:
		return snippet{
			tok.Content,
			tok.Params,
		}
	}
	return snippet{}
}

func readFromWeb(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func main() {

	var urlString string
	var fileString string
	var langString string

	flag.StringVar(&urlString, "url", "", `The url to a markdown file`)
	flag.StringVar(&fileString, "file", "", `The path to a markdown file`)
	flag.StringVar(&langString, "lang", "", `Select specific language to print`)
	flag.Parse()

	var readMe []byte

	switch {
	case urlString != "":
		var err error
		readMe, err = readFromWeb(urlString)
		if err != nil {
			log.Fatalf(err.Error())
		}

	case fileString != "":
		var err error
		readMe, err = ioutil.ReadFile(fileString)
		if err != nil {
			log.Fatalf(err.Error())
		}

	default:
		log.Fatalln("Please, provide url or file parameter.")
	}

	md := markdown.New(markdown.XHTMLOutput(true), markdown.Nofollow(true))
	tokens := md.Parse(readMe)

	for _, t := range tokens {
		snippet := getSnippet(t)

		if snippet.content != "" && (langString == "" || langString == snippet.lang) {
			fmt.Printf("##### Lang : %s ###### \n", snippet.lang)
			fmt.Println(snippet.content)
		}
	}

}
