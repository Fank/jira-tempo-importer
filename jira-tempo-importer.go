package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Printf("hello, world\n")

	auth := base64.StdEncoding.EncodeToString([]byte(os.Getenv("USERNAME") + ":" + os.Getenv("PASSWORD")))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	date := time.Now().Add(-(7 * 24 * time.Hour))
	url := fmt.Sprintf("https://jira.mcl.de/rest/tempo-timesheets/3/worklogs/?dateFrom=%4d-%02d-%02d&username=%s", date.Year(), date.Month(), date.Day(), os.Getenv("USERNAME"))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Basic "+auth)
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var worklogs []worklog
	json.NewDecoder(res.Body).Decode(&worklogs)

	for i := range worklogs {
		newTypeValue := ""
		switch strings.ToLower(worklogs[i].Comment) {

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

		if len(worklogs[i].WorkAttributeValues) > 0 && (worklogs[i].WorkAttributeValues[0].Value == newTypeValue || worklogs[i].WorkAttributeValues[0].Value == "refinement") {
			fmt.Printf("OK %s\n", worklogs[i].Comment)
		} else {
			fmt.Printf("NOT OK %s -> %s\n", worklogs[i].Comment, newTypeValue)

			newWorklog := worklogs[i]

			newValues := make([]worklogAttributeValue, 1)
			newValues[0] = worklogAttributeValue{
				Value: newTypeValue,
				WorkAttribute: worklogAttributeValueAttribute{
					ID: 1,
				},
				WorklogID: newWorklog.ID,
			}

			newWorklog.WorkAttributeValues = newValues

			jsonString, err := json.Marshal(newWorklog)
			if err == nil {
				// fmt.Println(string(jsonString))
				req, err = http.NewRequest("PUT", newWorklog.Self, bytes.NewBuffer(jsonString))
				req.Header.Set("Authorization", "Basic "+auth)
				req.Header.Set("Content-Type", "application/json;charset=UTF-8")
				if err == nil {
					res, err = client.Do(req)
					if err != nil {
						fmt.Printf("PUT failed: %s\n", err)
					} else if res.StatusCode != 200 {
						fmt.Printf("PUT failed: Status code %d\n", res.StatusCode)
					}
				} else {
					fmt.Printf("PUT prepare failed: %s\n", err)
				}
			} else {
				fmt.Printf("JSON failed: %s\n", err)
			}
		}
	}
}
