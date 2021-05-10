package main

import (
	"bytes"
	"context"
	_ "embed"
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

	//go:embed main.tmpl
	tmpl string
)

func fetchStatus() (*Status, error) {
	if os.Getenv("R_STATUS_URL") == "" {
		return nil, fmt.Errorf("the status URL to fetch must be defined as the environment variable R_STATUS_URL")
	}

	req, err := http.NewRequest("GET", os.Getenv("R_STATUS_URL"), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", os.Getenv("R_USER_AGENT"))

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
	if empty(os.Getenv("R_CLIENT_ID")) || empty(os.Getenv("R_CLIENT_SECRET")) || empty(os.Getenv("R_USERNAME")) || empty(os.Getenv("R_PASSWORD")) || empty(os.Getenv("R_SUBREDDIT")) || empty(os.Getenv("R_USER_AGENT")) {
		return nil, fmt.Errorf("one or more required environment variables were empty")
	}

	return reddit.NewClient(
		reddit.WithCredentials(os.Getenv("R_CLIENT_ID"), os.Getenv("R_CLIENT_SECRET"), os.Getenv("R_USERNAME"), os.Getenv("R_PASSWORD")),
		reddit.WithUserAgent(os.Getenv("R_USER_AGENT")),
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

	status := "[](#StatusStartMarker)" + statusText

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

func run() error {
	status, err := fetchStatus()
	if err != nil || status == nil {
		return err
	}

	remap := remapStatus(status)

	tmpl, err := template.New("template").Parse(tmpl)
	if err != nil {
		return err
	}

	writer := bytes.NewBuffer([]byte{})

	if err := tmpl.Execute(writer, remap); err != nil {
		return err
	}

	byts, err := ioutil.ReadAll(writer)
	if err != nil {
		return err
	}

	if err := updateSidebar(string(byts)); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
