package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/shurcooL/github_flavored_markdown"
	"github.com/shurcooL/github_flavored_markdown/gfmstyle"
	"github.com/urfave/cli"
)

type Page struct {
	Path     string
	Markdown []byte
	Html     []byte
}

type Log struct {
	Path       string
	ParentPath string
	Text       string
}

func convertMdToHtml(readme []byte) (*Page, error) {

	enrichedReadme := enrichMd(getReadmeAsMarkdown())

	// Make sure that there is at least one newline before our heading

	readmeAsHtml := github_flavored_markdown.Markdown(enrichedReadme)
	return &Page{Path: "FIXME", Markdown: enrichedReadme, Html: readmeAsHtml}, nil
}

/**
 * Adds various informations to the README.md of the container.
 */
func enrichMd(readme []byte) []byte {

	appendixTemplate, error := template.ParseFiles("appendix_template.md")
	var appendixBuffer bytes.Buffer
	error = appendixTemplate.Execute(&appendixBuffer, nil)

	if error != nil {
		fmt.Printf("Error while processing template: >%s<.", error)
		return readme
	}

	return append(readme, appendixBuffer.Bytes()...)
}

// `/var/log/AdminServer` - sshd log file
func parseLogFilePathsFromMd(pathToFile string) (*[]string, error) {

	var foundPath []string

	pageAsMarkdown, error := os.Open(pathToFile)

	if error != nil {
		return nil, error
	}

	matcher := regexp.MustCompile("`(.*)`.*-.*log.*")
	scanner := bufio.NewScanner(pageAsMarkdown)
	for scanner.Scan() {
		line := scanner.Text()
		match := matcher.FindStringSubmatch(line)
		if len(match) > 0 {
			foundPath = append(foundPath, match[1])
			fmt.Printf("Found path >%s< in file >%s<.\n", match[1], pathToFile)
		}
	}

	fmt.Printf("Found %d paths in file >%s<\n", len(foundPath), pathToFile)
	return &foundPath, nil
}

func readLogFiles(pathToLogFiles *[]string) (*[]Log, error) {

	var logs []Log

	for _, pathToLogFile := range *pathToLogFiles {
		file, err := os.Open(pathToLogFile)

		if err != nil {
			fmt.Printf("Error while reading log file >%s<: >%s<\n", pathToLogFile, err)
		} else {
			fmt.Printf("Successfully read log file >%s<.\n", pathToLogFile)
		}

		defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		log := Log{Path: pathToLogFile, ParentPath: pathToLogFile, Text: strings.Join(lines, "\n")}
		logs = append(logs, log)
	}

	return &logs, nil
}

func init() {
	// By default logger is set to write to stderr device.
	//log.SetOutput(os.Stdout)
}

func main() {

	app := cli.NewApp()

	app.Email = "christoph.zauner@NLLK.net"
	app.Author = "Christoph Zauner"
	app.Version = "0.1.0-alpha"
	// cntrinfod, cntinfod
	app.Usage = "Container Insight: HTTP daemon which exposes and augments the containers REAMDE.md"

	app.Action = func(c *cli.Context) error {

		fmt.Printf("Starting daemon on port 8080...\n")

		http.HandleFunc("/", handler)
		http.HandleFunc("/log", logHandler)

		// Serve the "/assets/gfm.css" file.
		http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(gfmstyle.Assets)))

		http.ListenAndServe(":8080", nil)

		return nil
	}

	app.Run(os.Args)
}

func getReadmeAsMarkdown() []byte {

	var possibleLocations [2]string = [2]string{"/README.md",
		"testdata/README.md"}

	var fd *os.File
	var err error

	for _, possibleLocation := range possibleLocations {
		fd, err = os.Open(possibleLocation)
		if err == nil {
			fmt.Printf("Found README.md file: >%s<\n", possibleLocation)
			break
		} else {
			fmt.Printf("Error when looking for README.md at >%s<: >%s<\n", possibleLocation, err)
		}
	}
	defer fd.Close()

	if err != nil {
		fmt.Printf("Could not find README.md file.\n")
		os.Exit(-1)
	}

	var readme []byte
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		readme = append(readme, scanner.Bytes()...)
	}

	return readme
}

func handler(w http.ResponseWriter, r *http.Request) {
	page, _ := convertMdToHtml(getReadmeAsMarkdown())

	io.WriteString(w, `<html><head><meta charset="utf-8"><link href="/assets/gfm.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
	w.Write(page.Html)
	io.WriteString(w, `</article></body></html>`)
}

func logHandler(w http.ResponseWriter, r *http.Request) {

	pathToLogFiles, _ := parseLogFilePathsFromMd("testdata/README.md")
	logs, _ := readLogFiles(pathToLogFiles)

	t, error := template.ParseFiles("log.html")
	error = t.Execute(w, logs)

	if error != nil {
		fmt.Printf("Error while processing template: >%s<.", error)
	}
}
