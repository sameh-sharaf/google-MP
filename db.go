package main

import (
	"os"

	"bitbucket.org/ascensionlab/internalhelpers"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func getGAClientIDs(accountID string) ([]string, error) {
	if accountID == "" {
		return nil, nil
	}

	req := &dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("TABLE_ACCOUNT_CLIENTS")),
		KeyConditionExpression: aws.String("accountId = :accountId"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":accountId": &dynamodb.AttributeValue{
				S: aws.String(accountID),
			},
		},
	}
	out, err := db.Query(req)
	if err != nil {
		return nil, err
	}
	if out.Items == nil {
		return nil, nil
	}

	sessionIDs := []string{}
	for _, item := range out.Items {
		sessionIDs = append(sessionIDs, internalhelpers.S(item["ga.clientId"]))
	}

	return sessionIDs, nil
}
