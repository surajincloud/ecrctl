package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/spf13/cobra"
	awspkg "github.com/surajincloud/ecrctl/pkg/awsutil"
)

// repoCmd represents the repo command
var untagRepoCmd = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"repo"},
	Short:   "untag an ECR repository",
	Long: `untag an ECR repository
	For example:
		ecrctl untag repo your-repo
		`,
	RunE: untagRepo,
}

func untagRepo(cmd *cobra.Command, args []string) error {
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

	tag, err := cmd.Flags().GetString("tag")
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

	_, err = client.UntagResource(ctx, &ecr.UntagResourceInput{
		ResourceArn: repoDetail.Repositories[0].RepositoryArn,
		TagKeys:     []string{tag},
	})
	if err != nil {
		return err
	}
	fmt.Println("repo untagged successfully.")

	return nil

}
func init() {
	untagCmd.AddCommand(untagRepoCmd)
	untagRepoCmd.PersistentFlags().String("region", "", "region")
	untagRepoCmd.PersistentFlags().String("profile", "", "profile")
	untagRepoCmd.PersistentFlags().String("repo", "", "repo")
	untagRepoCmd.PersistentFlags().StringP("tag", "t", "", "tags")

}
