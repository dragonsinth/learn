package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
)

const (
	awsBucket = "atolio-external-lab"
	awsPrefix = "20240209-sblum/"
	awsRegion = "us-west-2"
)

func NewAwsClient() (*AwsClient, error) {
	// Create a session using your credentials and region.
	awsSess, err := session.NewSession(&aws.Config{
		Region:     aws.String(awsRegion), // replace with your preferred region
		MaxRetries: aws.Int(5),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	return &AwsClient{session: awsSess, uploader: s3manager.NewUploader(awsSess)}, nil
}

type AwsClient struct {
	session  *session.Session
	uploader *s3manager.Uploader
}

func (ac *AwsClient) Upload(ctx context.Context, path string, data []byte) error {
	_, err := ac.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(awsBucket),
		Key:    aws.String(awsPrefix + path),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func (ac *AwsClient) CleanReset(ctx context.Context) error {
	s3Svc := s3.New(ac.session)

	resp, err := s3Svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(awsBucket),
		Prefix: aws.String(awsPrefix),
	})

	if err != nil {
		return fmt.Errorf("failed to list objects: %w", err)
	}

	// Loop through the objects and delete them one by one.
	for _, item := range resp.Contents {
		_, err = s3Svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(awsBucket),
			Key:    item.Key,
		})

		if err != nil {
			return fmt.Errorf("failed to delete object: %w", err)
		}

		log.Printf("Successfully deleted object: %s\n", *item.Key)
	}
	return nil
}
