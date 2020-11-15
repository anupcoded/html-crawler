package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func headingCount(heading []string) map[string]int {
	level := make(map[string]int)
	for _, num := range heading {
		level[num] = level[num] + 1
	}
	return level
}

func main() { //main function

	baseURL := "https://www.w3.org/TR/xhtml1/"
	var heading []string
	var internal, external int
	var form bool
	headingPattern, _ := regexp.Compile(`h[0-9]`)
	linkPattern, _ := regexp.Compile("http")

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
loop:
	for {
		tt := z.Next()
		t := z.Token()
		//fmt.Println(internal)
		switch {
		case tt == html.ErrorToken:
			break loop

		case tt == html.DoctypeToken:
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
			isTitle := t.Data == "title"
			if isTitle {
				z.Next()
				fmt.Println("Title:", string(z.Text()))
			}

			//headingPattern, _ := regexp.Compile(`h[0-9]`)
			isHeading := headingPattern.MatchString(t.Data)
			if isHeading {
				heading = append(heading, t.Data)
				//z.Next()
				//h1 := isHeading.FindAllStringIndex("A B C B A", -1)
				//fmt.Println("Heading:", heading)
				//return string(heading)
			}

			isAnchor := t.Data == "a"
			if isAnchor {

				//linkPattern, _ := regexp.Compile("http")
				for _, a := range t.Attr {
					if a.Key == "href" {
						isExternal := linkPattern.MatchString(a.Val)
						if isExternal {
							external++
						} else {
							internal++
						}
						//fmt.Println("Found href:", a.Val)
						break
					}
				}
			}

			isForm := t.Data == "form"
			if isForm {
				form = true
			}
		}
	}

	for headingLevel, count := range headingCount(heading) {
		fmt.Println(headingLevel, ":", count)
	}
	fmt.Println("Internal Links:", internal, " External Links:", external)
	fmt.Println("Form Found: ", form)
	//fmt.Println(heading)
	//response.Body.Close()

	/* a_page = get_url(a_page, "http://google.com")
	http.HandleFunc("/", index_handler)//function based on the path "/"
	http.ListenAndServe(":8000", nil)// (port,server which is nill) */
}
