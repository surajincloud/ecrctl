/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/spf13/cobra"
	awspkg "github.com/surajincloud/ecrctl/pkg/awsutil"
)

// repoCmd represents the repo command
var createRepoCmd = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"repo"},
	Short:   "Create an ECR repository",
	Long: `Create an ECR repository
	For example:
		ecrctl create repo your-repo
		`,
	RunE: createRepo,
}

func createRepo(cmd *cobra.Command, args []string) error {
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
	// Create an EKS client using the loaded configuration
	client := ecr.NewFromConfig(cfg)

	out, err := client.CreateRepository(ctx, &ecr.CreateRepositoryInput{
		RepositoryName: aws.String(repoName),
	})
	if err != nil {
		return err
	}
	fmt.Println("repo created successfully.")
	fmt.Println("repo ARN: ", aws.ToString(out.Repository.RepositoryArn))
	fmt.Println("repo URI: ", aws.ToString(out.Repository.RepositoryUri))

	return nil

}
func init() {
	createCmd.AddCommand(createRepoCmd)
	createRepoCmd.PersistentFlags().String("region", "", "region")
	createRepoCmd.PersistentFlags().String("profile", "", "profile")
	createRepoCmd.PersistentFlags().String("repo", "", "repo")

}
