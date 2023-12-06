package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var httpRegex = regexp.MustCompile("http://[-._%/[:alnum:]?:=+]+")
var httpsRegex = regexp.MustCompile("https://[-._%/[:alnum:]?:=+]+")

// config
var retryCount = 5

func checkURLLiveness(url string, retryCount int) error {
	for i := 0; i < retryCount; i++ {
		resp, err := http.Head(url)
		if err != nil {
			return err
		}
		if resp.StatusCode/100 == 2 {
			// ok
			return nil
		}
		log.Printf("code = %d, url = %s\n", resp.StatusCode, url)
		if i == retryCount-1 {
			return errors.New("invalid status code")
		} else {
			// exponential backoff
			time.Sleep((1 << i) * time.Second)
		}
	}
	return nil
}

func checkFile(path string) (err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	all := httpRegex.FindAll(content, -1)
	var livenessErrors uint64 = 0
	for _, v := range all {
		url := string(v)
		log.Printf("%s: HTTP link: url = %s\n", path, url)
		if thisError := checkURLLiveness(url, retryCount); thisError != nil {
			livenessErrors++
			log.Printf("%s: not alive: url = %s, thiserror = %v\n", path, url, thisError)
		}
	}

	all = httpsRegex.FindAll(content, -1)
	for _, v := range all {
		url := string(v)
		if thisError := checkURLLiveness(url, retryCount); thisError != nil {
			livenessErrors++
			log.Printf("%s: not alive: url = %s, thiserror = %v\n", path, url, thisError)
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
		log.Printf("git ls-files failed")
		os.Exit(2)
	}
	paths := strings.Split(strings.ReplaceAll(string(output), "\r\n", "\n"), "\n")

	for _, path := range paths[:len(paths)-1] { // excludes the last element after the last newline
		info, err := os.Stat(path)
		if err != nil {
			numErrors++
			log.Printf("path = %s, %v\n", path, err)
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
				log.Printf("%v\n", err)
			}
		}
		continue
	}
	if numErrors > 0 {
		os.Exit(1)
	}
}
