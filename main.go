package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/bugsnag/bugsnag-go"
)

type GoogleMP struct {
	AccountID      string `json:"accountId"`
	TagName        string `json:"tagName"`
	CampaignID     string `json:"campaignId"`
	OriginSource   string `json:"originSource"`
	OriginMedium   string `json:"originMedium"`
	OriginContent  string `json:"originContent"`
	OriginTerm     string `json:"originTerm"`
	OriginCampaign string `json:"originCampaign"`
}

func main() {
	for {
		wg := sync.WaitGroup{}
		wg.Add(2)
		go getMessages(getSQSQueueURL(os.Getenv("QUEUE_INPUT")))
		go process()
		wg.Wait()
	}
}

func process() {
	for {
		msg := <-sqsMessages

		var record GoogleMP
		err := json.Unmarshal([]byte(*msg.Body), &record)
		if err != nil {
			bugsnag.Notify(err)
			continue
		}

		gaClientIDs, err := getGAClientIDs(record.AccountID)
		if err != nil {
			bugsnag.Notify(err)
			continue
		}

		errors := 0
		for _, gaClientID := range gaClientIDs {
			err := notifyConversion(gaClientID, record)
			if err != nil {
				errors++
			} else {
				log.Printf("GA notified for account: %s - clientId: %s - campaign: %s - tag: %s \n", record.AccountID, gaClientID, record.CampaignID, record.TagName)
			}
		}

		if errors == 0 {
			deleteMessage(os.Getenv("QUEUE_INPUT"), msg)
		}
	}
}
