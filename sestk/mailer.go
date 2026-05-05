package sestk

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

type SesMailer interface {
	SendEmail(context.Context, *sesv2.SendEmailInput, ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error)
}

func NewMailer(ctx context.Context, region string) SesMailer {
	var client SesMailer
	if region == "log" {
		client = makeLogClient()
	} else {
		client = makeAwsClient(ctx, region)
	}

	return client
}

// Construct a simple mock email client that logs message params with log.Printf.
func makeLogClient() SesMailer {
	return &logClient{}
}

type logClient struct {
	messageIdCounter int
}

func (c *logClient) SendEmail(
	ctx context.Context,
	params *sesv2.SendEmailInput,
	optFns ...func(*sesv2.Options,
	)) (*sesv2.SendEmailOutput, error) {

	c.messageIdCounter++

	log.Printf("SendEmail(%+v)", params)

	return &sesv2.SendEmailOutput{
		MessageId: aws.String(fmt.Sprintf("mock_message_%06d", c.messageIdCounter)),
	}, nil
}

// Construct a real AWS SES client that can send emails for real.
func makeAwsClient(ctx context.Context, region string) SesMailer {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		func(o *config.LoadOptions) error {
			o.Region = region
			return nil
		},
	)
	if err != nil {
		log.Fatalf("unable to load config for aws: %v", err.Error())
	}

	return sesv2.NewFromConfig(cfg)
}
