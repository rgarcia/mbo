package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Cancel struct {
	VisitID *string
}

func (cmd *Cancel) Name() string { return "cancel" }
func (cmd *Cancel) DefineFlags(fs *flag.FlagSet) {
	cmd.VisitID = fs.String("visitid", "", "Visit ID to cancel.")
}
func (cmd *Cancel) Run() {
	log.SetOutput(LogOutput())

	if *cmd.VisitID == "" {
		fmt.Println("Must specify visitid.")
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

	resp, err := client.Get(fmt.Sprintf("%s/ASP/adm/adm_res_canc.asp?visitID=%s&cType=1", MBO_URL, *cmd.VisitID))
	if err != nil || resp.StatusCode != 200 {
		log.Println(err)
		log.Println(resp)
		fmt.Println("Error performing cancel.")
	}
	defer resp.Body.Close()

	fmt.Println("Cancelled visit.")
}
