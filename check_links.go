package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var httpRegex = regexp.MustCompile("http://[-._%/[:alnum:]]+")
var httpsRegex = regexp.MustCompile("https://[-._%/[:alnum:]]+")

func checkURLLiveness(url string) error {
	resp, err := http.Head(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "code = %d, url = %s\n", resp.StatusCode, url)
		return errors.New("invalid status code")
	}
	return nil
}

func checkFile(path string) (err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	all := httpRegex.FindAll(content, -1)
	var httpErrors uint64 = 0
	for _, v := range all {
		url := string(v)
		httpErrors++
		fmt.Fprintf(os.Stderr, "%s: HTTP link: url = %s\n", path, url)
	}
	if httpErrors > 0 {
		err = fmt.Errorf("detected HTTP links: path = %s, prev error = %w", path, err)
	}

	all = httpsRegex.FindAll(content, -1)
	var livenessErrors uint64 = 0
	for _, v := range all {
		url := string(v)
		if thisError := checkURLLiveness(url); thisError != nil {
			livenessErrors++
			fmt.Fprintf(os.Stderr, "%s: not alive: url = %s, thiserror = %v\n", path, url, thisError)
		}
	}
	if livenessErrors > 0 {
		err = fmt.Errorf("liveness check failed: path = %s, prev error = %w", path, err)
	}

	return err
}

func main() {
	// All text files' extensions
	extensions := []string{
		".c",
		".cpp",
		".go",
		".h",
		".java",
		".mod",
		".md",
		".py",
		".rs",
		".sh",
		".txt",
	}
	numErrors := 0
	cmd := exec.Command("git", "ls-files")
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "git ls-files failed")
		os.Exit(2)
	}
	paths := strings.Split(strings.ReplaceAll(string(output), "\r\n", "\n"), "\n")

	for _, path := range paths[:len(paths)-1] { // excludes the last element after the last newline
		info, err := os.Stat(path)
		if err != nil {
			numErrors++
			fmt.Fprintf(os.Stderr, "path = %s, %v\n", path, err)
			continue
		}
		if info.IsDir() {
			continue
		}
		ext := filepath.Ext(path)
		ok := false
		for _, e := range extensions {
			if ext == e {
				ok = true
				break
			}
		}
		if ok {
			if err := checkFile(path); err != nil {
				numErrors++
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		}
		continue
	}
	if numErrors > 0 {
		os.Exit(1)
	}
}
