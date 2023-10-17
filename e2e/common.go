package e2e

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mysteriumnetwork/feedback/constants"
)

type s3Downloader struct {
	downloader *manager.Downloader
	client     *s3.Client
	bucket     string
}

func newS3Downloader(bucket string) (*s3Downloader, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not load AWS configuration: %w", err)
	}

	s3client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String("http://localhost:9090")
		o.Region = constants.EuCentral1RegionID
	})
	downloader := &manager.Downloader{
		S3: s3client,
	}
	return &s3Downloader{
		downloader: downloader,
		client:     s3client,
		bucket:     bucket,
	}, nil
}

func (s3d *s3Downloader) getFileContent(t *testing.T, filename string) ([]byte, error) {
	paginator := s3.NewListObjectsV2Paginator(s3d.client, &s3.ListObjectsV2Input{
		Bucket: &s3d.bucket,
	})
	page, err := paginator.NextPage(context.Background())
	for ; err == nil; page, err = paginator.NextPage(context.Background()) {
		for _, obj := range page.Contents {
			if strings.Contains(*obj.Key, filename) {
				buf := manager.NewWriteAtBuffer([]byte{})
				_, err := s3d.downloader.Download(context.Background(), buf, &s3.GetObjectInput{
					Bucket: &s3d.bucket,
					Key:    obj.Key,
				})
				if err != nil {
					return nil, fmt.Errorf("download failed: %w", err)
				}
				return buf.Bytes(), nil
			}
		}
	}
	if err != nil {
		return nil, fmt.Errorf("pagination error: %w", err)
	}
	return nil, fmt.Errorf("file not found")
}
