package awsutil

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"k8s.io/apimachinery/pkg/util/duration"
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

func GetAge(creationTime time.Time) string {

	currentTime := time.Now()

	age := currentTime.Sub(creationTime)

	return duration.HumanDuration(age)
}
