package awsutil

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

type RepositoriesTags struct {
	Key, Value string
}

type Repositories struct {
	Name string
	Arn  string
	Uri  string
	Tags []RepositoriesTags
}

func ListRepos(ctx context.Context, client *ecr.Client) ([]Repositories, error) {
	var repoList []Repositories
	output, err := client.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{})
	if err != nil {
		return []Repositories{}, err
	}

	for _, i := range output.Repositories {

		repoTags, err := ListTags(ctx, client, *i.RepositoryArn)
		if err != nil {
			return []Repositories{}, err
		}

		repoList = append(repoList, Repositories{
			Name: aws.ToString(i.RepositoryName),
			Arn:  aws.ToString(i.RepositoryArn),
			Uri:  aws.ToString(i.RepositoryUri),
			Tags: repoTags,
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
