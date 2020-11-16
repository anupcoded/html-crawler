package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"text/template"

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

func fetchInfo(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("info.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		baseURL := strings.Join(r.Form["url"], " ")

		//baseURL := "https://www.w3.org/TR/xhtml1/"
		var heading []string
		var internal, external, inacessible int
		var form bool
		var version, title string

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
					var link string
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
							//fmt.Println("Found href:", link)
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

		fmt.Fprintf(w, "Title: %s\n", title)
		fmt.Fprintf(w, "HTML Version: %s\n", version)
		for headingLevel, count := range headingCount(heading) {
			fmt.Fprintf(w, "%s:%d\n", headingLevel, count)
		}
		fmt.Fprintf(w, "Internal Links: %d , External Links: %d \n", internal, external)
		fmt.Fprintf(w, "Inaccesible Links: %d\n", inacessible)
		fmt.Fprintf(w, "Login Form Found: %t\n", form)
	}
}

func main() { //main function

	http.HandleFunc("/", fetchInfo)
	http.ListenAndServe(":8000", nil)

	/* baseURL := "https://www.w3.org/TR/xhtml1/"
		var heading []string
		var internal, external, inacessible int
		var form bool

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

				headingPattern, _ := regexp.Compile(`h[0-9]`)
				isHeading := headingPattern.MatchString(t.Data)
				if isHeading {
					heading = append(heading, t.Data)
				}

				isAnchor := t.Data == "a"
				if isAnchor {

					linkPattern, _ := regexp.Compile("http|https")
					var link string
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
							//fmt.Println("Found href:", link)
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
		fmt.Println("Inaccesible Links:", inacessible)
		fmt.Println("Login Form Found: ", form) */
	//fmt.Println(heading)
	//response.Body.Close()

	/* http.HandleFunc("/", fetchInfo)
	http.ListenAndServe(":8000", nil) */
	/* a_page = get_url(a_page, "http://google.com")
	http.HandleFunc("/", index_handler)//function based on the path "/"
	http.ListenAndServe(":8000", nil)// (port,server which is nill) */
}
