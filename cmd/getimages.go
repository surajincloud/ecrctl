package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/spf13/cobra"
	awspkg "github.com/surajincloud/ecrctl/pkg/awsutil"
)

// imagesCmd represents the images command
var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Get ECR images",
	Long: `List ECR images
    Examples:

	# List all images with URI
	ecrctl get images --repo <repo-name>
    `,
	RunE: GetImages,
}

func GetImages(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	// read flag values
	region, err := cmd.Flags().GetString("region")
	if err != nil {
		fmt.Println("no region", err)
	}
	repo, err := cmd.Flags().GetString("repo")
	if err != nil {
		fmt.Println("no repo", err)
	}

	// aws config
	cfg, err := awspkg.GetAWSConfig(ctx, region)
	if err != nil {
		log.Fatal(err)
	}
	// Create an EKS client using the loaded configuration
	client := ecr.NewFromConfig(cfg)
	imageList, err := awspkg.ListImages(ctx, client, repo)
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 5, 2, 3, ' ', tabwriter.TabIndent)
	defer w.Flush()
	fmt.Fprintln(w, "TAG", "\t", "CRITICAL", "\t", "HIGH", "\t", "MEDIUM", "\t", "SIZE", "\t", "AGE")
	for _, i := range imageList {
		fmt.Fprintln(w, i.Tag, "\t", i.CriticalVulnerability, "\t", i.HighVulnerability, "\t", i.MediumVulnerability, "\t", i.Size, "\t", i.Age)
	}
	return nil
}
func init() {
	getCmd.AddCommand(imagesCmd)
	imagesCmd.PersistentFlags().String("region", "", "region")
	imagesCmd.PersistentFlags().String("repo", "", "repo")
}
