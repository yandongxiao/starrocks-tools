package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var _rootDir string

// Prompts:
// 我想使用 Golang 编写一个工具，它以目录为输入，遍历目录下的 Markdown 文件，从markdown文件中读取并检查链接，如果失效，就输出到标准输出。 请帮我写出这个工具。
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <dir>")
		os.Exit(1)
	}
	dir := os.Args[1]
	_rootDir = dir
	checkDir(dir)
}

// checkLink sends a HEAD request to the given url and returns true if the status code is 200, false otherwise.
func checkLink(file, url string) bool {
	if strings.Contains(url, "20") ||
		strings.Contains(url, "#") ||
		strings.Contains(url, "@") {
		// skip links like ../../sql-reference/sql-statements/data-manipulation/CREATE%20ROUTINE%20LOAD.md
		// skip links like ../../sql-reference/sql-statements/data-manipulation/CREATE%20ROUTINE%20LOAD.md#load-data
		// skip links like mailto:user@hello.iam.gserviceaccount.com
		return true
	}

	// if the url is a local file path, make sure the file exists
	if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		if strings.HasPrefix(url, "/") {
			url = filepath.Join(_rootDir, url)
		} else {
			dir := filepath.Dir(file)
			url = filepath.Join(dir, url)
		}

		// make sure the file represented by the url exists
		if _, err := os.Stat(url); err != nil {
			// try to append .md suffix to the url
			if !strings.HasSuffix(url, ".md") {
				url += ".md"
			}
			if _, err := os.Stat(url); err != nil {
				return false
			}
			return true
		}
		return true
	}

	// we can not send a HEAD request, because some servers do not support HEAD requests
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	return resp.StatusCode == 200
}

// extractLinks returns a slice of urls extracted from the given markdown content.
// It assumes that the urls are enclosed by parentheses and preceded by a square bracket.
func extractLinks(content string) []string {
	re := regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(content, -1)
	var links []string
	for _, match := range matches {
		links = append(links, match[1])
	}
	return links
}

// checkFile reads the content of the given markdown file and checks the links in it.
// It prints the file name and the broken links to the standard output.
func checkFile(file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	links := extractLinks(string(content))
	var broken []string
	for _, link := range links {
		if !checkLink(file, link) {
			broken = append(broken, link)
		}
	}
	if len(broken) > 0 {
		fmt.Printf("%s link s broken:\n", file)
		for _, link := range broken {
			fmt.Printf("\t%s\n", link)
		}
	}
}

// checkDir walks through the given directory and checks the markdown files in it.
func checkDir(dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".md" {
			checkFile(path)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
