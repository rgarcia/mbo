package main

import (
	"encoding/gob"
	"fmt"
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
