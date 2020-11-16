# html-crawler
GO Web application which takes a website URL as an input and provides general information about the contents of the page

- HTML Version
- Page Title
- Headings count by level
- Amount of internal and external links
- Amount of inaccessible links
- If a page contains a login form

## Getting Started

These instructions will get you to run the application on your local machine for testing purposes.

## Running the application
    
1. In the target folder run the app by using the cmd
    ```
    go run htcrawl.go
    ```

2. Allow access for the windows defender firewall prompt.

3. Navigate to the webpage ["http;//localhost:8000/"](http://localhost:8000/) and enter the url for the webpage for which the information is to be fetched.