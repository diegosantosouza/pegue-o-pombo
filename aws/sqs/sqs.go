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

func DeleteSQSMessage(client *sqs.Client, arn string, receiptHandle string) error {
	queueURL, err := getQueueURL(client, arn)
	if err != nil {
		return err
	}

	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err = client.DeleteMessage(context.Background(), input)
	if err != nil {
		log.Printf("Error deleting message from queue: %v", err)
		return err
	}

	log.Println("Message deleted successfully")
	return nil
}

// getQueueURL obtém a URL da fila SQS usando o ARN da fila.
func getQueueURL(client *sqs.Client, arn string) (string, error) {
	result, err := client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(arn), // Passe o ARN diretamente como QueueName
	})
	if err != nil {
		log.Printf("Error getting queue URL: %v", err)
		return "", err
	}
	return *result.QueueUrl, nil
}
