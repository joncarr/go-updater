package scraper

import (
	"os"
	"strings"
	"unicode"

	"github.com/anaskhan96/soup"
)

// GetRootNode returns the root node of the provided URL
func GetRootNode(url string) soup.Root {
	resp, err := soup.Get(url)
	if err != nil {
		os.Exit(1)
	}
	return soup.HTMLParse(resp)
}

// GetLatestVersionNumber returns the latest version number of Go
// as scraped from golang.org/dl
func GetLatestVersionNumber(root soup.Root) string {
	verNum := root.Find("h2", "class", "toggleButton")

	sb := strings.Builder{}

	for _, c := range verNum.Text() {
		if unicode.IsDigit(c) || c == '.' {
			sb.WriteRune(c)
		}
	}

	return sb.String()
}
