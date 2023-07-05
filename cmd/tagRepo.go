package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/spf13/cobra"
	awspkg "github.com/surajincloud/ecrctl/pkg/awsutil"
)

// repoCmd represents the repo command
var tagRepoCmd = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"repo"},
	Short:   "tag an ECR repository",
	Long: `tag an ECR repository
	For example:
		ecrctl tag repo your-repo
		`,
	RunE: tagRepo,
}

func tagRepo(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	// read flag values
	region, err := cmd.Flags().GetString("region")
	if err != nil {
		fmt.Println("no region", err)
	}

	profile, err := cmd.Flags().GetString("profile")
	if err != nil {
		fmt.Println("no profile", err)
	}

	tags, err := cmd.Flags().GetString("tag")
	if err != nil {
		fmt.Println("no tags", err)
	}

	if len(args) == 0 {
		fmt.Println("please pass repo name")
		os.Exit(1)

	}
	// repo Name
	repoName := args[0]

	// aws config
	cfg, err := awspkg.GetAWSConfig(ctx, region, profile)
	if err != nil {
		return err
	}
	// tag an EKS client using the loaded configuration
	client := ecr.NewFromConfig(cfg)

	repoDetail, err := client.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{
		RepositoryNames: []string{repoName},
	})
	if err != nil {
		return err
	}

	keyValue := strings.Split(tags, "=")

	_, err = client.TagResource(ctx, &ecr.TagResourceInput{
		ResourceArn: repoDetail.Repositories[0].RepositoryArn,
		Tags: []types.Tag{
			{
				Key:   &keyValue[0],
				Value: &keyValue[1],
			},
		},
	})
	if err != nil {
		return err
	}
	fmt.Println("repo tagged successfully.")

	return nil

}
func init() {
	tagCmd.AddCommand(tagRepoCmd)
	tagRepoCmd.PersistentFlags().String("region", "", "region")
	tagRepoCmd.PersistentFlags().String("profile", "", "profile")
	tagRepoCmd.PersistentFlags().String("repo", "", "repo")
	getCmd.PersistentFlags().StringP("tag", "t", "", "tags")

}
