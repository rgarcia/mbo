package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

type Register struct {
	ID   *string
	Date *string
}

func (cmd *Register) Name() string     { return "register" }
func (cmd *Register) Synopsis() string { return "Register for a class" }
func (cmd *Register) DefineFlags(fs *flag.FlagSet) {
	cmd.ID = fs.String("id", "", "Class ID")
	cmd.Date = fs.String("date", "", "Class date")
}
func (cmd *Register) Run() {
	log.SetOutput(LogOutput())

	if *cmd.ID == "" {
		fmt.Println("Must provide classid.")
		return
	}
	if *cmd.Date == "" {
		fmt.Println("Must provide class date.")
		return
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

	// Load the reservation page
	// This contains the full POST URL we need to hit
	preRegURL := fmt.Sprintf("%s/ASP/res_a.asp?classId=%s&classDate=%s",
		MBO_URL, *cmd.ID, *cmd.Date)
	log.Println("Getting", preRegURL)
	resp, err := client.Get(preRegURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	contentsb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	contents := string(contentsb)

	// search for call to submitResForm()--this contains full URL for the POST to register
	re := regexp.MustCompile("submitResForm\\('res_deb\\.asp\\?(.*)',")
	submatch := re.FindStringSubmatch(contents)
	log.Println("submatch", submatch)
	if len(submatch) == 0 {
		log.Println(re)
		log.Println(contents)
		fmt.Println("Session expired, please log in again.")
		return
	}
	regURL := fmt.Sprintf("%s/ASP/res_deb.asp?%s", MBO_URL, submatch[1])
	log.Println("posting to", regURL)
	resp, err = client.PostForm(regURL, url.Values{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println(resp)
		fmt.Println("Unknown error")
		return
	}
	fmt.Println("Successfully registered for class")
}
