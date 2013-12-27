package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
	"time"
)

type Schedule struct {
}

func (cmd *Schedule) Name() string { return "schedule" }
func (cmd *Schedule) DefineFlags(fs *flag.FlagSet) {
}
func (cmd *Schedule) Run() {
	log.SetOutput(LogOutput())

	// Load session
	mboSession, err := LoadMBOSession()
	if err != nil {
		fmt.Println(err)
		return
	}
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: cookieJar}
	mbo_url, _ := url.Parse(MBO_URL)
	client.Jar.SetCookies(mbo_url, mboSession.Cookies)

	// Get schedule
	resp, err := client.Get(fmt.Sprintf("%s/ASP/my_sch.asp", MBO_URL))
	if err != nil || resp.StatusCode != 200 {
		log.Println(err)
		log.Println(resp)
		fmt.Println("Error getting schedule")
	}
	defer resp.Body.Close()

	// Read into buffer so it can be read multiple times
	var body bytes.Buffer
	if _, err := body.ReadFrom(resp.Body); err != nil {
		fmt.Println(err)
		return
	}

	// Detect logged in welcome message
	if !IsLoggedIn(body.Bytes()) {
		fmt.Println("Session has expired; please log in again.")
		return
	}

	// Format in space-separated columns of minimal width 5 and at least
	// one blank of padding (so wider column entries do not touch each other)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 1, ' ', 0)
	header := []string{
		"Day",
		"Time",
		"Class",
		"Coach",
		"Visit ID",
	}
	fmt.Fprintln(w, strings.Join(header, "\t"))
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body.Bytes()))
	doc.Find("#mySchTable tr.scheduleRow").Each(func(i int, s *goquery.Selection) {
		day, err := time.Parse("Mon 01/02/2006", Strip(s.Find("td.dateCell").Text()))
		if err != nil {
			log.Println("Could not parse time from", Strip(s.Find("td.dateCell").Text()))
		}

		// Find Visit ID in cancel cell
		// <a href="javascript:jsConfirm(44270, 1)">Cancel</a>
		visitID := ""
		html, _ := s.Find("td.cancelCell").Html()
		re := regexp.MustCompile("jsConfirm\\((\\d+),")
		submatch := re.FindStringSubmatch(html)
		log.Println("submatch", submatch)
		if len(submatch) == 0 {
			log.Println("could not find visit id", re, html)
		} else {
			visitID = submatch[1]
		}
		fmt.Fprintln(w, strings.Join([]string{
			day.Format("Mon Jan 2"),
			Strip(s.Find("td.timeCell").Text()),
			Strip(s.Find("td.classNameCell").Text()),
			Strip(s.Find("td.teacherCell").Text()),
			visitID,
		}, "\t"))
	})
	w.Flush()
}
