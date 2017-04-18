package main

import (
	"net/http"
	"net/url"
	"os"

	"bitbucket.org/ascensionlab/internalhttp"
	"github.com/bugsnag/bugsnag-go"
)

func notifyConversion(clientID string, record GoogleMP) error {
	params := url.Values{}
	params.Set("v", "1")
	params.Set("tid", os.Getenv("TRACKING_ID"))
	params.Set("t", "event")
	params.Set("ds", "crm")
	params.Set("ea", "tag")
	params.Set("ec", record.CampaignID)
	params.Set("el", record.TagName)
	params.Set("cid", clientID)
	params.Set("ev", tagValues[record.TagName])
	params.Set("cn", record.OriginCampaign)
	params.Set("cs", record.OriginSource)
	params.Set("cm", record.OriginMedium)
	params.Set("ck", record.OriginTerm)
	params.Set("cc", record.OriginContent)
	params.Set("ci", record.CampaignID)

	req := &http.Request{}
	client, err := internalhttp.NewClient(req)
	if err != nil {
		bugsnag.Notify(err)
		return err
	}

	resp, err := client.PostForm("https://www.google-analytics.com/collect", params)
	if client.Error(resp, err) != nil {
		bugsnag.Notify(err, bugsnag.MetaData{
			"debug": {
				"client ID":   clientID,
				"campaign ID": record.CampaignID,
				"tag name":    record.TagName,
			},
		})
		return err
	}

	return nil
}
