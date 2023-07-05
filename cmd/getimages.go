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
var getImagesCmd = &cobra.Command{
	Use:     "images",
	Aliases: []string{"image"},
	Short:   "Get ECR images",
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

	profile, err := cmd.Flags().GetString("profile")
	if err != nil {
		fmt.Println("no profile", err)
	}

	repo, err := cmd.Flags().GetString("repo")
	if err != nil {
		fmt.Println("no repo", err)
	}

	// aws config
	cfg, err := awspkg.GetAWSConfig(ctx, region, profile)
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
	fmt.Fprintln(w, "TAG", "\t", "VULNERABILITIES", "\t", "SIZE", "\t", "AGE")
	for _, i := range imageList {
		fmt.Fprintln(w, i.Tag, "\t", fmt.Sprintf("C: %d, H: %d, M: %d", i.CriticalVulnerability, i.HighVulnerability, i.MediumVulnerability), "\t", i.Size, "\t", i.Age)
	}
	return nil
}
func init() {
	getCmd.AddCommand(getImagesCmd)
	getImagesCmd.PersistentFlags().String("region", "", "region")
	getImagesCmd.PersistentFlags().String("profile", "", "profile")
	getImagesCmd.PersistentFlags().String("repo", "", "repo")
}
