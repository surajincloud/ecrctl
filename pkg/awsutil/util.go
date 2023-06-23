package awsutil

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func GetAWSConfig(ctx context.Context, region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return aws.Config{}, err
	}
	if cfg.Region == "" {
		// get region
		region, err = GetRegion(region)
		if err != nil {
			return aws.Config{}, err
		}
		cfg.Region = region
	}
	return cfg, nil
}

func GetRegion(region string) (string, error) {
	if region == "" {
		region = os.Getenv("AWS_REGION")
		if region == "" {
			return "", fmt.Errorf("please pass region name with --region or with AWS_REGION environment variable")
		}
		return region, nil
	}
	return region, nil
}
