package sqs

import (
	"context"
	"fmt"
	"log"
	"regexp"

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
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(arnToURL(arn)),
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

func arnToURL(arn string) string {
	// Padrão para extrair informações do ARN
	arnPattern := regexp.MustCompile(`arn:aws:sqs:([^:]+):(\d+):(.+)`)
	matches := arnPattern.FindStringSubmatch(arn)

	if len(matches) != 4 {
		return "" // Não foi possível extrair informações necessárias
	}

	region := matches[1]
	accountID := matches[2]
	queueName := matches[3]

	// Constrói a URL da fila SQS com base nas informações extraídas
	queueURL := fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", region, accountID, queueName)

	return queueURL
}
