package main

import "io/ioutil"
import "net/http"
import "html/template"

import "github.com/russross/blackfriday"

type Page struct {
	Path string
	Markdown []byte
	Html []byte
}

type Log struct {
	Path string
	RawContent string
	ModifiedContent string
}

func convertMdToHtml(pathToFile string) (*Page, error)  {
	
	pageAsMarkdown, error := ioutil.ReadFile(pathToFile)
	pageAsHtml := blackfriday.MarkdownBasic(pageAsMarkdown)
	
	if error != nil {
		return nil, error
	}

	return &Page{Path: pathToFile, Markdown: pageAsMarkdown, Html: pageAsHtml}, nil
}

func enrichMd(pathToFile string) (*Page, error)  {
	
	pageAsMarkdown, error := ioutil.ReadFile(pathToFile)
	pageAsHtml := blackfriday.MarkdownBasic(pageAsMarkdown)
	
	if error != nil {
		return nil, error
	}

	return &Page{Path: pathToFile, Markdown: pageAsMarkdown, Html: pageAsHtml}, nil
}

// `/var/log/AdminServer` --- sshd log file
func parseLogFilePathsFromMd(pathToFile string) (*Log, error)  {
	
	pageAsMarkdown, error := ioutil.ReadFile(pathToFile)
	pageAsHtml := blackfriday.MarkdownBasic(pageAsMarkdown)
	
	if error != nil {
		return nil, error
	}

	return &Page{Path: pathToFile, Markdown: pageAsMarkdown, Html: pageAsHtml}, nil
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/log", logHandler)
	http.ListenAndServe(":8080", nil)	
}


func handler(w http.ResponseWriter, r *http.Request) {
	page, _ := convertMdToHtml("README.md")
	w.Write(page.Html)
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	
	

	t, _ := template.ParseFiles("log.html")
	t.Execute(w, p)	
}
