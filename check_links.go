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

	"github.com/BurntSushi/toml"
)

const configFilePath = "./check_links_config.toml"

var httpRegex = regexp.MustCompile("http://[-._%/[:alnum:]?:=+]+")
var httpsRegex = regexp.MustCompile("https://[-._%/[:alnum:]?:=+]+")

type Config struct {
	RetryCount int `toml:"retry_count"`
	// All text files' extensions
	TextFileExtensions []string `toml:"text_file_extensions"`
	Ignores            []Ignore `toml:"ignores"`
}

type Ignore struct {
	URL                    string   `toml:"url"`
	Codes                  []int    `toml:"codes"`
	Reason                 string   `toml:"reason"`
	ConsideredAlternatives []string `toml:"considered_alternatives"`
}

func (c *Config) Validate() error {
	if len(c.TextFileExtensions) == 0 {
		return errors.New("text_file_extensions cannot be empty")
	}
	for _, ignore := range c.Ignores {
		if ignore.URL == "" {
			return errors.New("url cannot be empty")
		}
		if len(ignore.Codes) == 0 {
			return errors.New("codes cannot be empty")
		}
		if ignore.Reason == "" {
			return errors.New("reason cannot be empty")
		}
		if len(ignore.ConsideredAlternatives) == 0 {
			return errors.New("considered_alternatives cannot be empty")
		}
	}
	return nil
}

func readConfig(configFilePath string) (*Config, error) {
	var config Config
	bytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	toml.Decode(string(bytes), &config)
	return &config, nil
}

// if ignore != nil, ignore.Codes will be used instead of the 2xx criterion.
func checkURLLiveness(url string, retryCount int, ignore *Ignore) error {
	for i := 0; i < retryCount; i++ {
		resp, err := http.Head(url)
		if err != nil {
			return err
		}
		if ignore != nil {
			ok := false
			for _, code := range ignore.Codes {
				if resp.StatusCode == code {
					ok = true
					break
				}
			}
			if ok {
				// ok, but because ignore != nil, we need a log
				log.Printf("ok: code = %d, url = %s, ignore = %v\n", resp.StatusCode, url, ignore)
				return nil
			}
		} else {
			if resp.StatusCode/100 == 2 {
				// ok
				return nil
			}
		}
		log.Printf("code = %d, url = %s, ignore = %v\n", resp.StatusCode, url, ignore)
		if i == retryCount-1 {
			return errors.New("invalid status code")
		} else {
			// exponential backoff
			time.Sleep((1 << i) * time.Second)
		}
	}
	return nil
}

func checkFile(path string, retryCount int, ignores map[string]*Ignore) (err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	all := httpRegex.FindAll(content, -1)
	var livenessErrors uint64 = 0
	for _, v := range all {
		url := string(v)
		ignore := ignores[url]
		log.Printf("%s: HTTP link: url = %s\n", path, url)
		if thisError := checkURLLiveness(url, retryCount, ignore); thisError != nil {
			livenessErrors++
			log.Printf("%s: not alive: url = %s, thiserror = %v\n", path, url, thisError)
		}
	}

	all = httpsRegex.FindAll(content, -1)
	for _, v := range all {
		url := string(v)
		ignore := ignores[url]
		if thisError := checkURLLiveness(url, retryCount, ignore); thisError != nil {
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
	config, err := readConfig(configFilePath)
	if err != nil {
		panic(err)
	}
	if err := config.Validate(); err != nil {
		panic(err)
	}
	ignores := make(map[string]*Ignore)
	for _, ignore := range config.Ignores {
		// For handling of https://go.dev/blog/loopvar-preview
		ignoreCopied := ignore
		ignores[ignore.URL] = &ignoreCopied
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
		for _, e := range config.TextFileExtensions {
			if ext == e {
				ok = true
				break
			}
		}
		if ok {
			if err := checkFile(path, config.RetryCount, ignores); err != nil {
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
