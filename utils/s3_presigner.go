package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"time"
)

type ArchiveResourceRequest struct {
	ArchiveID string
	ArchiveSlug string
	File string
	Download bool
}
// Downloads an item from an S3 Bucket
func PresignArchiveResource(request *ArchiveResourceRequest) string {
	config := GetConfig()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AwsRegion)},
	)

	// Create S3 service client
	svc := s3.New(sess)

	input := &s3.GetObjectInput{
		Bucket: aws.String(config.BucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", request.ArchiveID, request.File)),
	}

	if request.Download == true {
		filename := fmt.Sprintf("%s-%s", request.ArchiveSlug, request.File)
		input.ResponseContentDisposition = aws.String(fmt.Sprintf("attachment; filename=%s", filename))
	}

	req, _ := svc.GetObjectRequest(input)
	urlStr, err := req.Presign(30 * time.Minute)

	if err != nil {
		log.Println("Failed to sign request", err)
		return ""
	}

	return urlStr
}
