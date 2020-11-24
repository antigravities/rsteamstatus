package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"

	_ "github.com/joho/godotenv/autoload" // yes, I'm lazy
	"github.com/vartanbeno/go-reddit/reddit"
)

var translated map[string]string = map[string]string{
	"online":    "Online on Steam",
	"cms":       "Client",
	"webapi":    "Web API",
	"store":     "Store",
	"community": "Community",
}

var (
	httpClient *http.Client    = &http.Client{}
	ctx        context.Context = context.Background()
)

func fetchStatus() (*Status, error) {
	if os.Getenv("R_STATUS_URL") == "" {
		return nil, fmt.Errorf("The status URL to fetch must be defined as the environment variable R_STATUS_URL")
	}

	httpClient := &http.Client{}

	req, err := http.NewRequest("GET", os.Getenv("R_STATUS_URL"), nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_0_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	byts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	status := &Status{}
	err = json.Unmarshal(byts, status)
	if err != nil {
		return nil, err
	}

	return status, nil
}

func isGood(status string) string {
	if strings.ToLower(status) == "ok" || strings.ToLower(status) == "normal" || strings.HasSuffix(strings.ToLower(status), "million") || strings.HasPrefix(status, "9") || strings.HasPrefix(status, "100") {
		return "good"
	}

	return "bad"
}

func empty(item string) bool {
	return len(item) == 0
}

func makeReddit() (*reddit.Client, error) {
	if empty(os.Getenv("R_CLIENT_ID")) || empty(os.Getenv("R_CLIENT_SECRET")) || empty(os.Getenv("R_USERNAME")) || empty(os.Getenv("R_PASSWORD")) || empty(os.Getenv("R_SUBREDDIT")) {
		return nil, fmt.Errorf("One or more required environment variables were empty")
	}

	return reddit.NewClient(
		reddit.WithCredentials(os.Getenv("R_CLIENT_ID"), os.Getenv("R_CLIENT_SECRET"), os.Getenv("R_USERNAME"), os.Getenv("R_PASSWORD")),
		reddit.WithUserAgent("Golang:get.cutie.cafe/rsteamstatus:1.0 (by /u/antigravities)"),
	)
}

func updateSidebar(statusText string) error {
	bot, err := makeReddit()
	if err != nil {
		return err
	}

	page, _, err := bot.Wiki.Page(ctx, os.Getenv("R_SUBREDDIT"), "config/sidebar")
	if err != nil {
		return err
	}

	rangeX := strings.Index(page.Content, "[](#StatusStartMarker)")
	rangeY := strings.Index(page.Content, "[](#StatusEndMarker)")

	status := "[](#StatusStartMarker)" + statusText + "[](#StatusEndMarker)"

	content := strings.Replace(page.Content, page.Content[rangeX:rangeY], status, 1)

	_, err = bot.Wiki.Edit(ctx, &reddit.WikiPageEditRequest{
		Subreddit: os.Getenv("R_SUBREDDIT"),
		Page:      "config/sidebar",
		Content:   content,
		Reason:    "/r/Steam status update",
	})

	return err
}

func remapStatus(status *Status) *Remap {
	remap := &Remap{}
	// fuck it, ship it
	remap.Statuses = make(map[string]struct {
		Name   string
		Good   string
		Status string
	})

	for _, svc := range status.Services {
		remap.Statuses[svc[0].(string)] = struct {
			Name   string
			Good   string
			Status string
		}{
			Name:   translated[svc[0].(string)],
			Good:   isGood(svc[2].(string)),
			Status: svc[2].(string),
		}
	}

	return remap
}

func openTemplate(path string) (*template.Template, error) {
	byts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return template.New("template").Parse(string(byts))
}

func main() {
	status, err := fetchStatus()
	if err != nil || status == nil {
		panic(err)
	}

	remap := remapStatus(status)

	tmpl, err := openTemplate("main.tmpl")
	if err != nil {
		panic(err)
	}

	writer := bytes.NewBuffer([]byte{})

	if err := tmpl.Execute(writer, remap); err != nil {
		panic(err)
	}

	byts, err := ioutil.ReadAll(writer)
	if err != nil {
		panic(err)
	}

	if err := updateSidebar(string(byts)); err != nil {
		panic(err)
	}
}
