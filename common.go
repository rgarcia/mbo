package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// MBO_URL is the base url for the MINDBODY Online client portal.
const MBO_URL = "https://clients.mindbodyonline.com"

// MBOSession stores information needed to scrape MBO's website as a logged-in user.
type MBOSession struct {
	Cookies  []*http.Cookie
	StudioID string
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
