package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/spf13/cobra"
	awspkg "github.com/surajincloud/ecrctl/pkg/awsutil"
)

// repositoriesCmd represents the repositories command
var describeRepositoriesCmd = &cobra.Command{
	Use:     "repositories",
	Aliases: []string{"repo"},
	Short:   "Describe a repositories",
	Long: `Describe a repositories
    For example:

	ecrctl describe repositories <repo-name>
    `,
	RunE: describeRepo,
}

func describeRepo(cmd *cobra.Command, args []string) error {

	ctx := context.Background()
	// read flag values
	region, err := cmd.Flags().GetString("region")
	if err != nil {
		fmt.Println("no region", err)
	}

	if len(args) == 0 {
		fmt.Println("please pass repo name")
		os.Exit(1)

	}
	// repo Name
	repoName := args[0]

	// aws config
	cfg, err := awspkg.GetAWSConfig(ctx, region)
	if err != nil {
		return err
	}
	// Create an EKS client using the loaded configuration
	client := ecr.NewFromConfig(cfg)

	repo, err := awspkg.DescribeRepository(ctx, client, repoName)
	if err != nil {
		return err
	}
	fmt.Println("Name:           ", repo.Name)
	fmt.Println("URI:            ", repo.Uri)
	fmt.Println("ARN:            ", repo.Arn)
	fmt.Println("TAG Mutability: ", repo.TagMutability)
	fmt.Println("Created At:     ", repo.CreatedAt)
	fmt.Println("Tags:           ", TagsToString(repo.Tags))
	fmt.Println("Scan On Push:   ", repo.ScanOnPush)

	return nil
}

func TagsToString(repoTagList []awspkg.RepositoriesTags) string {
	var tagList []string
	for _, k := range repoTagList {
		tagList = append(tagList, fmt.Sprintf("%s=%s", k.Key, k.Value))
	}

	return strings.Join(tagList, ",")
}

func init() {
	describeCmd.AddCommand(describeRepositoriesCmd)
	describeCmd.PersistentFlags().String("region", "", "region")
}
