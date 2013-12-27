package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// MBO_URL is the base url for the MINDBODY Online client portal.
const MBO_URL = "https://clients.mindbodyonline.com"

// MBOSession stores information needed to scrape MBO's website as a logged-in user.
type MBOSession struct {
	Cookies  []*http.Cookie
	StudioID string
}

func LoadMBOSession() (mboSession MBOSession, err error) {
	file, err := os.Open(fmt.Sprintf("%s/.mindbodyonline", os.Getenv("HOME")))
	if err != nil {
		return mboSession, fmt.Errorf("Must be logged in.")
	}
	defer file.Close()
	dec := gob.NewDecoder(file)
	err = dec.Decode(&mboSession)
	if err != nil {
		return mboSession, fmt.Errorf("Must be logged in.")
	}
	return mboSession, nil
}

// LogOutput determines where we should send logs (if anywhere).
func LogOutput() (logOutput io.Writer) {
	logOutput = ioutil.Discard
	if os.Getenv("MBO_LOG") != "" {
		logOutput = os.Stderr
		if logPath := os.Getenv("MBO_LOG_PATH"); logPath != "" {
			var err error
			logOutput, err = os.Create(logPath)
			if err != nil {
				panic(err)
			}
		}
	}
	return
}

// IsLoggedIn detects if the user is logged in. Pass it any html on a logged in page.
func IsLoggedIn(html []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(html))
	selection := doc.Find("#top-wel-sp")
	if selection.Length() != 1 {
		return false
	}
	return true
}

// Strip strips the result of goquery.Text(), removing leading/trailing whitespace and nbsp's
func Strip(str string) string {
	return strings.Replace(strings.TrimSpace(str), " ", " ", -1)
}
