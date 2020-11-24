package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"
)

var translated map[string]string = map[string]string{
	"online":    "Online on Steam",
	"cms":       "Client",
	"webapi":    "Web API",
	"store":     "Store",
	"community": "Community",
}

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

func remapStatus(status *Status) *Remap {
	remap := &Remap{}
	remap.Statuses = make(map[string]struct {
		Name   string
		Good   string
		Status string
	})

	for _, svc := range status.Services {
		// fuck it, ship it
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

	fmt.Println(string(byts))
}
