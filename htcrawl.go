package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

//Page represents a
/* type Page struct {
	Url           string
	Version       string
	Title         string
	HeadingLevel  uint16
	InternalLinks uint16
	ExternalLinks uint16
	InaccLinks    uint16
	Login         bool
} */

/* type HttpError struct {
	Errorinfo string
} */

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//Parse will take an HTML url and returns information
/* func Parse(resp *http.Response, depth int) []Page {
	pageData := html.NewTokenizer(resp.Body)
	links := []Page{}

	var start *html.Token
	var text string

	for{
		_=pageData.Next()
		token := page.Token()
		if token.Type == html.ErrorToken {
			break
		}

		if start != nil && token.Type == html.TextToken {
			text = fmt.Sprintf("%s%s", text, token.Data)
		}

		if token.dataAtom == atom.A {
			switch token.type {
			case html.StartTagToken:
				if len(token.Attr) > 0 {
					start = &token
				}
			case html.EndTagToken:
				if start == nil {
					log.Warnf("Link end found without start:%s",text)
					continue
				}
				link := NewLink(*start, text, depth)
				if link.Valid(){
					links = append(links, link)
					log.Deugf("link found %v", link)
				}

				start = nil
				text = ""
			}
		}
	}

	log.Debug(links)
	return links
}

func NewLink(tag html.Token, text string, depth int) Link {
	link := link
} */

/* func urlParser(Page p, url string) Page {
	p.page_url = url
	return p
} */

/* func get_info(url string) {
	resp, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(resp.Body)

	//b := resp.Body
	defer b.Close()

	tokenizer := html.NewTokenizer(resp.Body)

	for {
		toktyp:= tokenizer.Next()
		if toktyp == html.ErrorToken {
			if tokenizer.Err() != io.EOF {
				if tokenizer.Err() != io.EOF {
        			WARN.Println(fmt.Sprintf(“HTML error found in %s due to “,
                    currentUrl, tokenizer.Err()))
				}
				return
			}
		}
		token := tokenizer.token()
		switch  toktyp {

		}

	}

	string_body := string(bytes)
	fmt.Println(string_body)
	resp.Body.Close()

}  */

func main() { //main function

	baseURL := "https://www.w3.org/TR/xhtml1/"

	//skip https
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: config,
	}

	netClient := &http.Client{
		Transport: transport,
	}

	response, err := netClient.Get(baseURL)

	checkErr(err)

	z := html.NewTokenizer(response.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.DoctypeToken:
			t := z.Token()
			fmt.Println(t.Data)
			switch {
			case t.Data == "html":
				fmt.Println("HTML Version: HTML5")
			case t.Data == "HTML":
				fmt.Println("HTML Version: HTML5 and beyond")
			case strings.Contains(t.Data, "HTML 4.01"):
				fmt.Println("HTML Version: HTML 4.01")
			case strings.Contains(t.Data, "HTML 3.2"):
				fmt.Println("HTML Version: HTML 3.2")
			case strings.Contains(t.Data, "HTML 2.0"):
				fmt.Println("HTML Version: HTML 2.0")
			case strings.Contains(t.Data, "XHTML 1.0"):
				fmt.Println("HTML Version: XHTML 1.0")
			case strings.Contains(t.Data, "XHTML 1.1"):
				fmt.Println("HTML Version: XHTML 1.1")
			}

		case tt == html.StartTagToken:
			t := z.Token()

			isTitle := t.Data == "title"
			if isTitle {
				z.Next()
				fmt.Println("Title:", string(z.Text()))
			}

			isAnchor := t.Data == "a"
			if isAnchor {

				for _, a := range t.Attr {
					if a.Key == "href" {
						fmt.Println("Found href:", a.Val)
						break
					}
				}
			}

			isForm := t.Data == "form"
			if isForm {
				fmt.Println("form found")
			}

		}

	}

	//response.Body.Close()

	/* a_page = get_url(a_page, "http://google.com")
	http.HandleFunc("/", index_handler)//function based on the path "/"
	http.ListenAndServe(":8000", nil)// (port,server which is nill) */

}
