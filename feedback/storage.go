/*
 * Copyright (C) 2019 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package feedback

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	"github.com/gofrs/uuid"
)

// Storage file storage
type Storage struct {
	uploader *s3manager.Uploader
	bucket   string
}

// NewStorageOpts options to initialize Storage
type NewStorageOpts struct {
	EndpointURL string
	Bucket      string
}

// New creates a new Storage
func New(opts *NewStorageOpts) (storage *Storage, err error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load AWS configuration: %w", err)
	}
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(opts.EndpointURL)
	cfg.Region = endpoints.EuCentral1RegionID
	s3client := s3.New(cfg)
	s3client.ForcePathStyle = true

	s3uploader := &s3manager.Uploader{
		S3:                s3client,
		PartSize:          s3manager.DefaultUploadPartSize,
		Concurrency:       s3manager.DefaultUploadConcurrency,
		LeavePartsOnError: false,
		MaxUploadParts:    s3manager.MaxUploadParts,
	}

	return &Storage{
		uploader: s3uploader,
		bucket:   opts.Bucket,
	}, nil
}

// Upload uploads file to storage and returns its URL
func (s *Storage) Upload(filepath string) (url *url.URL, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not open file for reading %q: %w", filepath, err)
	}

	fileKey, err := s.generateRemoteKey(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not generate remote key for upload: %w", err)
	}

	result, err := s.uploader.Upload(&s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileKey),
	})

	if err != nil {
		return nil, fmt.Errorf("could not upload file %q: %w", filepath, err)
	}
	uploadURL, err := url.Parse(result.Location)
	if err != nil {
		return nil, fmt.Errorf("could not resolve uploaded file URL: %w", err)
	}

	return uploadURL, nil
}

func (*Storage) generateRemoteKey(filepath string) (remoteKey string, err error) {
	randomId, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	remoteKey = fmt.Sprintf("%s-%s-%s",
		time.Now().Format("2006-01-02"),
		randomId.String(),
		path.Base(filepath),
	)
	return remoteKey, nil
}
