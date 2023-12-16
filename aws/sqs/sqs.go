package sqs

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NewSQSSession() (*sqs.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return sqs.NewFromConfig(cfg), nil
}

func DeleteSQSMessage(client *sqs.Client, queueURL string, receiptHandle string) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := client.DeleteMessage(context.Background(), input)
	if err != nil {
		log.Printf("Error deleting message from queue: %v", err)
		return err
	}

	log.Println("Message deleted successfully")
	return nil
}
