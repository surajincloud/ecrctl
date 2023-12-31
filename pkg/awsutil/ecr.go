package awsutil

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/docker/go-units"
)

func DescribeImage(ctx context.Context, client *ecr.Client, imageName string, repoName string) (Image, error) {
	imageList, err := ListImages(ctx, client, repoName)
	if err != nil {
		return Image{}, err
	}

	for _, i := range imageList {
		if i.Tag == imageName {
			scans, err := client.DescribeImageScanFindings(ctx, &ecr.DescribeImageScanFindingsInput{
				RepositoryName: &repoName,
				ImageId: &types.ImageIdentifier{
					ImageTag: &imageName,
				},
			})
			if err != nil {
				fmt.Println("error while scan findings")
			}
			i.BasicScanFindings = scans.ImageScanFindings.Findings
			i.EnhancedScanFindings = scans.ImageScanFindings.EnhancedFindings
			return i, nil
		}
	}
	return Image{}, nil
}

func ListImages(ctx context.Context, client *ecr.Client, repoName string) ([]Image, error) {

	describeImagesOutput, err := client.DescribeImages(ctx, &ecr.DescribeImagesInput{
		RepositoryName: aws.String(repoName),
	})

	if err != nil {
		return []Image{}, err
	}

	// Sort the images in descending order based on the image push timestamp
	sort.Slice(describeImagesOutput.ImageDetails, func(i, j int) bool {
		return describeImagesOutput.ImageDetails[i].ImagePushedAt.After(*describeImagesOutput.ImageDetails[j].ImagePushedAt)
	})

	var imageList []Image

	for k, i := range describeImagesOutput.ImageDetails {

		if k <= 4 {

			var tags string
			if len(i.ImageTags) == 0 {
				tags = "untagged"
			} else {
				tags = strings.Join(i.ImageTags, ",")
			}

			size := units.HumanSize(float64(*i.ImageSizeInBytes))
			imageList = append(imageList, Image{
				ArtifactMediaType:     *i.ArtifactMediaType,
				Digest:                *i.ImageDigest,
				Tag:                   tags,
				CriticalVulnerability: i.ImageScanFindingsSummary.FindingSeverityCounts["CRITICAL"],
				HighVulnerability:     i.ImageScanFindingsSummary.FindingSeverityCounts["HIGH"],
				MediumVulnerability:   i.ImageScanFindingsSummary.FindingSeverityCounts["MEDIUM"],
				Age:                   GetAge(*i.ImagePushedAt),
				Size:                  size,
				ScanStatus:            i.ImageScanStatus.Status,
				ScanStatusDesc:        *i.ImageScanStatus.Description,
			})
		}
	}
	return imageList, nil
}

func DescribeRepository(ctx context.Context, client *ecr.Client, repoName string) (Repository, error) {
	repoList, err := ListRepos(ctx, client)
	if err != nil {
		return Repository{}, err
	}
	for _, i := range repoList {
		if i.Name == repoName {
			return i, nil
		}
	}
	return Repository{}, nil
}

func ListRepos(ctx context.Context, client *ecr.Client) ([]Repository, error) {
	var repoList []Repository
	output, err := client.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{})
	if err != nil {
		return []Repository{}, err
	}

	for _, i := range output.Repositories {

		repoTags, err := ListTags(ctx, client, *i.RepositoryArn)
		if err != nil {
			return []Repository{}, err
		}
		repoList = append(repoList, Repository{
			Name:             aws.ToString(i.RepositoryName),
			Arn:              aws.ToString(i.RepositoryArn),
			Uri:              aws.ToString(i.RepositoryUri),
			Age:              GetAge(*i.CreatedAt),
			CreatedAt:        i.CreatedAt,
			Tags:             repoTags,
			TagMutability:    i.ImageTagMutability,
			EncryptionType:   i.EncryptionConfiguration.EncryptionType,
			EncryptionKMSKey: aws.ToString(i.EncryptionConfiguration.KmsKey),
			ScanOnPush:       i.ImageScanningConfiguration.ScanOnPush,
		})

	}

	return repoList, nil
}

func ListTags(ctx context.Context, client *ecr.Client, repoArn string) (repoTags []RepositoriesTags, err error) {
	out, err := client.ListTagsForResource(ctx, &ecr.ListTagsForResourceInput{
		ResourceArn: &repoArn,
	})
	if err != nil {
		return []RepositoriesTags{}, err
	}

	for _, k := range out.Tags {
		repoTags = append(repoTags, RepositoriesTags{
			Key:   aws.ToString(k.Key),
			Value: aws.ToString(k.Value),
		})
	}
	return repoTags, nil

}
