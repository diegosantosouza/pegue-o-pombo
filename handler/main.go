package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	sesCli "github.com/diegosantosouza/pegue-o-pombo/aws/ses"
	sqsCli "github.com/diegosantosouza/pegue-o-pombo/aws/sqs"
)

func ProcessMessage(ctx context.Context, sqsEvent events.SQSEvent) error {
	// Config SES client
	sesClient, err := sesCli.NewSESSession()
	if err != nil {
		log.Printf("Error to load sesClient: %v", err)
		return err
	}

	// Config SQS client
	sqsClient, err := sqsCli.NewSQSSession()
	if err != nil {
		log.Printf("Error to load sqsClient: %v", err)
		return err
	}

	for _, record := range sqsEvent.Records {
		var emailEvent sesCli.EmailEvent
		err := json.Unmarshal([]byte(record.Body), &emailEvent)
		if err != nil {
			log.Printf("Error decoding message SQS: %v", err)
			return err
		}

		log.Printf("Message received: %s", emailEvent)

		// Send email using SES
		err = sesCli.SendEmail(sesClient, emailEvent)
		if err != nil {
			log.Printf("Error sending email: %v", err)
			return err
		}

		// Delete message
		queueURL := record.EventSourceARN
		receiptHandle := record.ReceiptHandle
		err = sqsCli.DeleteSQSMessage(sqsClient, queueURL, receiptHandle)
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
