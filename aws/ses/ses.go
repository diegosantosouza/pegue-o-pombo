package ses

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

// EmailEvent SQS type
type EmailEvent struct {
	Subject          string   `json:"Subject"`
	Body             string   `json:"Body"`
	ToAddresses      []string `json:"ToAddresses"`
	FromEmailAddress string   `json:"FromEmailAddress"`
}

func NewSESSession() (*sesv2.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return sesv2.NewFromConfig(cfg), nil
}

func SendEmail(client *sesv2.Client, emailEvent EmailEvent) error {
	emailInput := &sesv2.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: emailEvent.ToAddresses,
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Data: aws.String(emailEvent.Body),
					},
				},
				Subject: &types.Content{
					Data: aws.String(emailEvent.Subject),
				},
			},
		},
		FromEmailAddress: aws.String(emailEvent.FromEmailAddress),
	}

	_, err := client.SendEmail(context.Background(), emailInput)
	if err != nil {
		return err
	}

	return nil
}
