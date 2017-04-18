package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/bugsnag/bugsnag-go"
)

func getSQSQueueURL(queueName string) string {
	params := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}
	resp, err := mq.GetQueueUrl(params)
	if err != nil {
		bugsnag.Notify(err, bugsnag.MetaData{
			"debug": {
				"Queue": queueName,
			},
		})
		return ""
	}

	return *resp.QueueUrl
}

func getMessages(queueURL string) {
	defer wg.Done()
	for {
		resp, err := mq.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: aws.Int64(1),
			WaitTimeSeconds:     aws.Int64(20),
		})
		if err != nil {
			bugsnag.Notify(err)
			time.Sleep(5 * time.Minute)
		}

		for _, message := range resp.Messages {
			sqsMessages <- message
			changeMessageVisibility(queueURL, message, 120)
		}
	}
}

func deleteMessage(queueName string, msg *sqs.Message) {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(getSQSQueueURL(queueName)),
		ReceiptHandle: msg.ReceiptHandle,
	}

	_, err := mq.DeleteMessage(params)
	if err != nil {
		bugsnag.Notify(err)
	}
}

func changeMessageVisibility(queueURL string, msg *sqs.Message, timeout int64) {
	params := &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(queueURL),
		ReceiptHandle:     aws.String(*msg.ReceiptHandle),
		VisibilityTimeout: aws.Int64(timeout),
	}
	_, err := mq.ChangeMessageVisibility(params)

	if err != nil {
		bugsnag.Notify(err)
		return
	}
}
