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

// repositoryCmd represents the repository command
var deleteRepoCmd = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"repo"},
	Short:   "Delete an ECR repository",
	Long: `Delete an ECR repository
	For example:
		ecrctl delete repo your-repo
		`,
	RunE: deleteRepo,
}

func deleteRepo(cmd *cobra.Command, args []string) error {
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

	_, err = client.DeleteRepository(ctx, &ecr.DeleteRepositoryInput{
		RepositoryName: aws.String(repoName),
		Force:          false,
	})

	if err != nil {
		fmt.Println("It seems that repository may contains some images.")
		fmt.Println("you will need --force flag to force the repository deletion, it will also delete the images from the repository.")
	} else {
		fmt.Println("repo deleted successfully.")
	}

	return nil
}
func init() {
	deleteCmd.AddCommand(deleteRepoCmd)
	deleteRepoCmd.PersistentFlags().String("region", "", "region")
	deleteRepoCmd.PersistentFlags().String("profile", "", "profile")
	deleteRepoCmd.PersistentFlags().String("repo", "", "repo")

}
