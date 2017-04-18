package main

import (
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/bugsnag/bugsnag-go"
)

var mq *sqs.SQS
var wg sync.WaitGroup
var db *dynamodb.DynamoDB
var sqsMessages chan *sqs.Message
var tagValues map[string]string

func init() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       "7adadaf776bb815dc3eb94585e71935a",
		ReleaseStage: os.Getenv("RELEASE_STAGE"),
	})

	if os.Getenv("QUEUE_INPUT") == "" {
		log.Fatalln("QUEUE_INPUT not set")
	}
	if os.Getenv("TABLE_ACCOUNT_CLIENTS") == "" {
		log.Fatalln("TABLE_ACCOUNT_CLIENTS not set")
	}
	if os.Getenv("TRACKING_ID") == "" {
		log.Fatalln("TRACKING_ID not set")
	}

	mq = sqs.New(session.New())
	db = dynamodb.New(session.New())
	sqsMessages = make(chan *sqs.Message)
	tagValues = map[string]string{
		"qualified": "1",
		"verified":  "10",
		"submitted": "25",
		"approved":  "100",
	}
}
