package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"text/template"

	"golang.org/x/net/html"
)

// Page represents few properties of a url.
type Page struct {
	Heading                         []string
	Internal, External, Inacessible int
	Form                            bool
	Version, Title                  string
}

//check err
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

//headingCount counts headings  based on levels
func headingCount(heading []string) map[string]int {
	level := make(map[string]int)
	for _, num := range heading {
		level[num] = level[num] + 1
	}
	return level
}

// getInfo for a particular URL
func getInfo(baseURL string) *Page {

	var heading []string
	var internal, external, inacessible int
	var form bool
	var version, title, link string

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

		switch {
		case tt == html.ErrorToken:
			break loop

		case tt == html.DoctypeToken:
			switch {
			case t.Data == "html":
				version = "HTML5"
			case t.Data == "HTML":
				version = "HTML5 and beyond"
			case strings.Contains(t.Data, "HTML 4.01"):
				version = "HTML 4.01"
			case strings.Contains(t.Data, "HTML 3.2"):
				version = "HTML 3.2"
			case strings.Contains(t.Data, "HTML 2.0"):
				version = "HTML 2.0"
			case strings.Contains(t.Data, "XHTML 1.0"):
				version = "XHTML 1.0"
			case strings.Contains(t.Data, "XHTML 1.1"):
				version = "XHTML 1.1"
			default:
				version = "Other"
			}

		case tt == html.StartTagToken:
			isTitle := t.Data == "title"
			if isTitle {
				z.Next()
				title = string(z.Text())
			}

			headingPattern, _ := regexp.Compile(`h[0-9]`)
			isHeading := headingPattern.MatchString(t.Data)
			if isHeading {
				heading = append(heading, t.Data)
			}

			isAnchor := t.Data == "a"
			if isAnchor {
				linkPattern, _ := regexp.Compile("http|https")
				for _, a := range t.Attr {
					if a.Key == "href" {
						isExternal := linkPattern.MatchString(a.Val)
						if isExternal {
							external++
							link = a.Val
						} else {
							internal++
							link = fmt.Sprintf("%s%s", baseURL, a.Val)
						}
						_, err := url.ParseRequestURI(link)
						if err != nil {
							inacessible++
						}
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

	return &Page{Heading: heading, Internal: internal, External: external, Inacessible: inacessible, Form: form, Version: version, Title: title}
}

//fetchURL fetches and gets information of URL
func fetchURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)

	if r.Method == "GET" {
		t, err := template.ParseFiles("info.gtpl")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		//execute template
		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		r.ParseForm()

		baseURL := strings.Join(r.Form["url"], " ")
		p := getInfo(baseURL)

		fmt.Fprintf(w, `<h1>HTML Page Information</h1>
		<p>Title: %s</p>
		<p>HTML Version: %s</p>
		<p>Headings count by level: </p>`, p.Title, p.Version)

		for headingLevel, count := range headingCount(p.Heading) {
			fmt.Fprintf(w, `%s: %d</p>`, headingLevel, count)
		}

		fmt.Fprintf(w, `<p>Internal Links: %d , External Links: %d </p>
		<p>Inaccesible Links: %d</p>
		<p>Login Form Found: %t</p>`, p.Internal, p.External, p.Inacessible, p.Form)
	}
}

func main() {

	http.HandleFunc("/", fetchURL)
	err := http.ListenAndServe(":8000", nil)
	checkErr(err)
}
