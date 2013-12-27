package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

type Ls struct {
	Date *string
}

func (cmd *Ls) Name() string { return "ls" }
func (cmd *Ls) DefineFlags(fs *flag.FlagSet) {
	cmd.Date = fs.String("date", "", "list classes as of this date. Format is MM/DD/YYYY. Default is today.")
}
func (cmd *Ls) Run() {
	log.SetOutput(LogOutput())

	if *cmd.Date == "" {
		timestr := time.Now().Format("02/01/2006")
		cmd.Date = &timestr
	}

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

	// List classes
	resp, err := client.Get(fmt.Sprintf("%s/ASP/main_class.asp?%s", MBO_URL, url.Values{
		"date": []string{*cmd.Date},
	}.Encode()))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Detect logged in welcome message
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	selection := doc.Find("#top-wel-sp")
	if selection.Length() != 1 {
		fmt.Println("Session has expired; please log in again.")
		return
	}

	// Format in space-separated columns of minimal width 5 and at least
	// one blank of padding (so wider column entries do not touch each other)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 1, ' ', 0)
	header := []string{
		"Day",
		"Start Time",
		"Sign Up",
		"Class ID",
		"Class Name",
		"Trainer Name",
		"Assistant",
		"Duration",
	}
	fmt.Fprintln(w, strings.Join(header, "\t"))

	// Strip leading/trailing whitespace and nbsp's
	strip := func(str string) string {
		return strings.Replace(strings.TrimSpace(str), " ", " ", -1)
	}

	var day *time.Time
	doc.Find("#classSchedule-mainTable tr").Each(func(i int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		if class != "evenRow" && class != "oddRow" {
			str := strip(s.Text())
			log.Printf("Header text: '%s'", str)
			time, err := time.Parse("Mon January 02, 2006", str)
			if err == nil {
				day = &time
			} else {
				log.Println("Could not parse time from", str)
			}
			return
		}
		if day == nil {
			// haven't parsed a date header in the table yet
			html, _ := s.Html()
			log.Println("Error, nil day", i, class, html)
			return
		}

		// Find a signup button--this contains class ID
		// <input type="button" name="but169" class="SignupButton" ...
		signup := s.Find("input.SignupButton")
		classID := ""
		if signup.Length() != 0 {
			classID, _ = signup.Attr("name")
			classID = classID[3:]
		}
		fmt.Fprintln(w, strings.Join([]string{
			day.Format("Mon Jan 2"),
			strip(s.Find("td:nth-child(1)").Text()),
			strip(s.Find("td:nth-child(2)").Text()),
			classID,
			strip(s.Find("td:nth-child(3)").Text()),
			strip(s.Find("td:nth-child(4)").Text()),
			strip(s.Find("td:nth-child(5)").Text()),
			strip(s.Find("td:nth-child(6)").Text()),
		}, "\t"))
	})
	w.Flush()

	return
}
