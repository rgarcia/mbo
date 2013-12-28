package main

import (
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Login logs into MINDBODY Online.
type Login struct {
	Username *string
	Password *string
	StudioID *string
}

func (cmd *Login) Name() string     { return "login" }
func (cmd *Login) Synopsis() string { return "Start a session with MBO" }
func (cmd *Login) DefineFlags(fs *flag.FlagSet) {
	cmd.Username = fs.String("u", "", "Username. Will prompt if not passed.")
	cmd.Password = fs.String("p", "", "Password. Will prompt if not passed.")
	cmd.StudioID = fs.String("studio", "", "Studio ID. Will prompt if not passed.")
}
func (cmd *Login) Run() {
	// Request credentials if not passed as arguments
	username := *cmd.Username
	if username == "" {
		fmt.Print("username: ")
		_, err := fmt.Scanln(&username)
		if err != nil {
			panic(err)
		}
	}

	password := *cmd.Password
	if password == "" {
		fmt.Print("password: ")
		_, err := fmt.Scanln(&password)
		if err != nil {
			panic(err)
		}
	}

	studioID := *cmd.StudioID
	if studioID == "" {
		fmt.Print("studio id: ")
		_, err := fmt.Scanln(&studioID)
		if err != nil {
			panic(err)
		}
	}

	// three steps to a successful login:
	// 1. browse to https://clients.mindbodyonline.com/ASP/ws.asp?studioId=x
	// 2. accept redirect to https://clients.mindbodyonline.com/ASP/ws.asp?studioId=x&sessionChecked=true
	// 3. post login form to https://clients.mindbodyonline.com/ASP/login_p.asp
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: cookieJar}
	mbo_url, _ := url.Parse(MBO_URL)
	_, err := client.Get(fmt.Sprintf("%s/ASP/ws.asp?%s", MBO_URL,
		url.Values{"studioId": []string{studioID}}.Encode()))
	if err != nil {
		panic(err)
	}
	_, err = client.Get(fmt.Sprintf("%s/ASP/ws.asp?%s", MBO_URL,
		url.Values{"studioId": []string{studioID}, "sessionChecked": []string{"true"}}.Encode()))
	cookieJar.SetCookies(mbo_url, []*http.Cookie{&http.Cookie{Name: "f5_cspm", Value: "1234"}})
	resp, err := client.PostForm(fmt.Sprintf("%s/ASP/login_p.asp?%s", MBO_URL,
		url.Values{
			"isLibAsync":        []string{"true"},
			"isJson":            []string{"true"},
			"libAsyncTimeStamp": []string{strconv.FormatInt(time.Now().Unix()*1000, 10)},
			"studioId":          []string{studioID},
		}.Encode()),
		url.Values{
			"requiredtxtUserName": []string{username},
			"requiredtxtPassword": []string{password},
		})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&data); err != nil {
		panic(err)
	}
	if _, ok := data["json"]; !ok {
		fmt.Println("Invalid credentials")
		return
	}
	if success, ok := data["json"].(map[string]interface{})["success"]; !ok || !(success.(bool)) {
		fmt.Println("Invalid credentials")
		return
	}
	fmt.Println("successfully logged in!")

	// Save session for later
	mboSession := MBOSession{
		StudioID: studioID,
		Cookies:  cookieJar.Cookies(mbo_url),
	}
	file, err := os.Create(fmt.Sprintf("%s/.mindbodyonline", os.Getenv("HOME")))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	err = enc.Encode(mboSession)
	if err != nil {
		panic(err)
	}

	return
}
