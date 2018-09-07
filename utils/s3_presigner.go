package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

// Downloads an item from an S3 Bucket
func PresignArchiveResource(archiveID uuid.UUID, file string) string {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	// Create S3 service client
	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("teyitlink"),
		Key:    aws.String(fmt.Sprintf("%s/%s", archiveID.String(), file)),
	})
	urlStr, err := req.Presign(30 * time.Minute)

	if err != nil {
		log.Println("Failed to sign request", err)
		return ""
	}

	return urlStr
}
