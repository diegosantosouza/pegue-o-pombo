package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// EmailEvent SQS type
type EmailEvent struct {
	Subject          string   `json:"Subject"`
	Body             string   `json:"Body"`
	ToAddresses      []string `json:"ToAddresses"`
	FromEmailAddress string   `json:"FromEmailAddress"`
}

func ProcessMessage(ctx context.Context, sqsEvent events.SQSEvent) error {
	// Config SES client
	sesClient, err := newSESSession()
	if err != nil {
		log.Printf("Error to load sesClient: %v", err)
		return err
	}

	// Config SQS client
	sqsClient, err := newSQSSession()
	if err != nil {
		log.Printf("Error to load sqsClient: %v", err)
		return err
	}

	for _, record := range sqsEvent.Records {
		var emailEvent EmailEvent
		err := json.Unmarshal([]byte(record.Body), &emailEvent)
		if err != nil {
			log.Printf("Error decoding message SQS: %v", err)
			return err
		}

		log.Printf("Message received: %s", emailEvent)

		// Send email using SES
		err = sendEmail(sesClient, emailEvent)
		if err != nil {
			log.Printf("Error sending email: %v", err)
			return err
		}

		// Delete message
		queueURL := record.EventSourceARN
		receiptHandle := record.ReceiptHandle
		err = deleteSQSMessage(sqsClient, queueURL, receiptHandle)
		if err != nil {
			log.Printf("Error deleting SQS message: %v", err)
			return err
		}
	}

	return nil
}

func main() {
	lambda.Start(ProcessMessage)
}
