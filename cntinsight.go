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
	"strconv"
	"strings"
	"time"

	"github.com/zaunerc/cntinsight/docker"
	"github.com/zaunerc/cntinsight/system"

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

	location := locateReadme()
	enrichedReadme := enrichMd(getReadmeAsMarkdown(location))

	readmeAsHtml := github_flavored_markdown.Markdown(enrichedReadme)
	return &Page{Path: location, Markdown: enrichedReadme, Html: readmeAsHtml}, nil
}

/**
 * Adds various informations to the README.md of the container.
 */
func enrichMd(readme []byte) []byte {

	appendixTemplate, error := template.ParseFiles("appendix_template.md")
	var appendixBuffer bytes.Buffer

	vars := map[string]interface{}{
		"ContainerHostname": system.FetchContainerHostname(),
		"HostHostname":      docker.FetchHostHostname(),
		"TcpSocketInfo":     system.FetchTcp46SocketInfo(),
		"UdpSocketInfo":     system.FetchUdp46SocketInfo(),
		"ProcessInfo":       system.FetchProcessInfo(),
		"ProcessTree":       system.FetchProcessTree(),
	}

	error = appendixTemplate.Execute(&appendixBuffer, vars)

	if error != nil {
		fmt.Printf("Error while processing template: >%s<.", error)
		return readme
	}

	// Make sure that there is at least one newline before our heading
	readme = append(readme, "\n"...)
	return append(readme, appendixBuffer.Bytes()...)
}

// Use getReadmeAsMarkdown
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

	var httpPort int

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "httpPort, p",
			Value:       2020,
			Usage:       "Listen on port `PORT` for HTTP connections",
			Destination: &httpPort,
		},
	}

	app.Email = "christoph.zauner@NLLK.net"
	app.Author = "Christoph Zauner"
	app.Version = "0.1.0-alpha"
	// cntrinfod, cntinfod
	app.Usage = "Container Insight: HTTP daemon which exposes and augments the containers REAMDE.md"

	app.Action = func(c *cli.Context) error {

		fmt.Printf("Starting HTTP daemon on port %d...\n", httpPort)

		http.HandleFunc("/", protect(handler))
		http.HandleFunc("/log", protect(logHandler))
		http.HandleFunc("/hostinfo", protect(hostInfoHandler))
		http.HandleFunc("/markdown", protect(markdownHandler))

		// Serve the "/assets/gfm.css" file.
		http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(gfmstyle.Assets)))

		http.ListenAndServe(":"+strconv.Itoa(httpPort), nil)

		return nil
	}

	app.Run(os.Args)
}

func protect(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if checkAuth(w, r) {
			fn(w, r)
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="Linux PAM"`)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
		}

	}
}

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	user, pass, _ := r.BasicAuth()
	fmt.Printf("User >%s< trying to authenticate\n", user)
	return user == "user" && pass == "pass"
}

func locateReadme() string {

	var possibleLocations [2]string = [2]string{"/README.md",
		"testdata/README.md"}

	var fd *os.File
	defer fd.Close()
	var err error

	var location string
	for _, possibleLocation := range possibleLocations {
		fd, err = os.Open(possibleLocation)
		if err == nil {
			fmt.Printf("Found README.md file: >%s<\n", possibleLocation)
			location = possibleLocation
			break
		} else {
			fmt.Printf("Error when looking for README.md at >%s<: >%s<\n", possibleLocation, err)
		}
	}

	if err != nil {
		fmt.Printf("Could not find README.md file.\n")
		os.Exit(-1)
	}

	return location
}

func getReadmeAsMarkdown(path string) []byte {

	var fd *os.File
	defer fd.Close()
	var err error

	fd, err = os.Open(path)
	if err != nil {
		fmt.Printf("Error while trying to open REAMDE.md file: %s\n", err)
		os.Exit(-1)
	}

	var readme []byte
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		readme = append(readme, scanner.Bytes()...)
		readme = append(readme, "\n"...)
	}

	return readme
}

func handler(w http.ResponseWriter, r *http.Request) {
	page, _ := convertMdToHtml(getReadmeAsMarkdown(locateReadme()))

	io.WriteString(w, `<html><head><meta charset="utf-8"><link href="/assets/gfm.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
	w.Write(page.Html)
	io.WriteString(w, `</article></body></html>`)
}

func logHandler(w http.ResponseWriter, r *http.Request) {

	pathToLogFiles, _ := parseLogFilePathsFromMd(locateReadme())
	logs, _ := readLogFiles(pathToLogFiles)

	t, error := template.ParseFiles("log.html")

	currentDateAndTime := time.Now().Format("2006-01-02 15:04:05")

	vars := map[string]interface{}{
		"Logs":               logs,
		"CurrentDateAndTime": currentDateAndTime,
	}

	error = t.Execute(w, vars)

	if error != nil {
		fmt.Printf("Error while processing template: >%s<.", error)
	}
}

func hostInfoHandler(w http.ResponseWriter, r *http.Request) {
	t, error := template.ParseFiles("hostinfo.html")

	hostinfo := docker.FetchHostInfo()
	error = t.Execute(w, hostinfo)

	if error != nil {
		fmt.Printf("Error while processing template: >%s<.", error)
	}
}

func markdownHandler(w http.ResponseWriter, r *http.Request) {

	readmeAsMarkdown := getReadmeAsMarkdown(locateReadme())
	enrichedReadme := enrichMd(readmeAsMarkdown)

	t, error := template.ParseFiles("markdown.html")
	error = t.Execute(w, string(enrichedReadme))
	if error != nil {
		fmt.Printf("Error while processing template: >%s<.", error)
	}
}
