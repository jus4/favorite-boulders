package utils

import (
  "bytes"
  "context"
  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/config"
  "github.com/aws/aws-sdk-go-v2/service/s3"
  "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

type AwsS3Uploader struct {
  Client *s3.Client
}

func NewAwsS3Uploader(ctx context.Context)(*AwsS3Uploader, error) {
  cfg, err := config.LoadDefaultConfig(context.TODO())
  if err != nil {
    return nil, err
  }

  client := s3.NewFromConfig(cfg)

  return &AwsS3Uploader{client}, nil
}

func (a *AwsS3Uploader) UploadTopoImage(ctx context.Context, key string, file []byte)(string, error) {
  uploader := manager.NewUploader(a.Client)
  result, uploadError := uploader.Upload(context.TODO(), &s3.PutObjectInput{
      Bucket: aws.String("suomi-topot-images"),
      Key:    aws.String("topo_images/"+key),
      Body:   bytes.NewReader(file),
  })

  if uploadError != nil {
    return "Failed to upload image", uploadError
  }

  return result.Location, nil
}
