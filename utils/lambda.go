package utils

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/satori/go.uuid"
)

type archiveRequest struct {
	ArchiveId  uuid.UUID `json:"archive_id"`
	RequestUrl string    `json:"request_url"`
}

type ArchiveResponsePayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func RunArchiveLambda(archiveId uuid.UUID, requestUrl string) (*ArchiveResponsePayload, error) {
	config := GetConfig()

	// Create Lambda service client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambda.New(sess, &aws.Config{Region: aws.String(config.AwsRegion)})

	// Create the archive request for Lambda
	request := archiveRequest{archiveId, requestUrl}
	payload, err := json.Marshal(request)

	// Error marshalling request
	if err != nil {
		return nil, err
	}

	result, err := client.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String(config.WorkerLambdaName),
		Payload:      payload,
	})

	// Error running the lambda
	if err != nil {
		return nil, err
	}

	// If the status code is NOT 200, the archiving failed
	if *result.StatusCode != 200 {
		return nil, errors.New("archiving failed")
	}

	var responsePayload ArchiveResponsePayload

	err = json.Unmarshal(result.Payload, &responsePayload)

	// Error unmarshaling response payload
	if err != nil {
		return nil, err
	}

	return &responsePayload, nil
}
