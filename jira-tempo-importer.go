package main

import (
	"crypto/tls"
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"bytes"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func main() {
	initFlags()
	err := updateJira()
	if err != nil {
		panic(err)
	}
}

var flagToken string
var flagUrl string
var flagUsername string

func initFlags() {
	flag.StringVar(&flagToken, "token", "", "Tempo access token")
	flag.StringVar(&flagUrl, "url", "", "URL to access the jira server")
	flag.StringVar(&flagUsername, "username", "", "User which logs should be updated")
	flag.Parse()
}

func updateJira() error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	date := time.Now().Add(-(7 * 24 * time.Hour))
	url := fmt.Sprintf(
		"%s/plugins/servlet/tempo-getWorklog/?format=xml&tempoApiToken=%s&userName=%s&dateFrom=%4d-%02d-%02d",
		flagUrl,
		flagToken,
		flagUsername,
		date.Year(), date.Month(), date.Day(),
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	//b, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	return err
	//}
	//logrus.Println(string(b))

	var worklogs worklogsXML
	err = xml.NewDecoder(res.Body).Decode(&worklogs)
	if err != nil {
		return err
	}

	for _, v := range worklogs.Worklogs {
		if v.BillingAttributes != "" {
			continue
		}

		newTypeValue := ""
		switch strings.ToLower(v.WorkDescription) {

		case "merge approval":
			fallthrough
		case "ma":
			newTypeValue = "codereview"

		case "support":
			fallthrough
		case "s":
			newTypeValue = "support"
		case "c":
			newTypeValue = "creation"

		case "bugfixes":
			fallthrough
		case "merge fixes":
			fallthrough
		case "code review":
			fallthrough
		case "cr":
			newTypeValue = "codereview"

		case "cl":
			newTypeValue = "clearing"
		case "r":
			newTypeValue = "refinement"

		default:
			fallthrough
		case "":
			newTypeValue = "development"
		}

		fmt.Printf("NOT OK %s -> %s\n", v.WorkDescription, newTypeValue)

		updates := worklogUpdatesXML{
			WorklogUpdates: []worklogUpdateXML{
				{
					ID:                v.ID,
					BillingAttributes: fmt.Sprintf("type=%s", newTypeValue),
					HashValue:         v.HashValue,
				},
			},
		}
		changes, err := xml.Marshal(updates)
		if err != nil {
			return err
		}
		logrus.Println(string(changes))

		url = fmt.Sprintf(
			"%s/plugins/servlet/tempo-updateWorklog/?tempoApiToken=%s",
			flagUrl,
			flagToken,
		)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(changes))
		if err != nil {
			return err
		}
		res, err := client.Do(req)
		if err != nil {
			return err
		}

		logrus.Warnln(res.Status)

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		logrus.Warnln(string(b))

		res.Body.Close()
	}

	return nil
}
