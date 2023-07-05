/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/spf13/cobra"
	awspkg "github.com/surajincloud/ecrctl/pkg/awsutil"
)

// imagesCmd represents the images command
var describeImagesCmd = &cobra.Command{
	Use:     "images",
	Aliases: []string{"image"},
	Short:   "describe an image",
	Long: `describe an image and get more information
	For example:
		ecrctl describe images your-image --repo your-repo
		`,
	RunE: describeImages,
}

func describeImages(cmd *cobra.Command, args []string) error {
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

	repoName, err := cmd.Flags().GetString("repo")
	if err != nil {
		fmt.Println("no repo", err)
	}

	if len(args) == 0 {
		fmt.Println("please pass repo name")
		os.Exit(1)

	}
	// image Name
	imageName := args[0]

	// aws config
	cfg, err := awspkg.GetAWSConfig(ctx, region, profile)
	if err != nil {
		return err
	}
	// Create an EKS client using the loaded configuration
	client := ecr.NewFromConfig(cfg)

	image, err := awspkg.DescribeImage(ctx, client, imageName, repoName)
	if err != nil {
		return err
	}
	fmt.Println("Name: ", imageName)
	fmt.Println("Repo: ", repoName)
	fmt.Println("Type: ", image.ArtifactMediaType)
	fmt.Println("Vulnerabilities:")
	w := tabwriter.NewWriter(os.Stdout, 5, 2, 3, ' ', tabwriter.TabIndent)
	defer w.Flush()
	fmt.Fprintln(w, "TAG", "\t", "SEVERITY", "\t", "URI")
	for _, i := range image.BasicScanFindings {
		fmt.Fprintln(w, *i.Name, "\t", i.Severity, "\t", *i.Uri)

	}
	return nil
}

func init() {
	describeCmd.AddCommand(describeImagesCmd)
	describeImagesCmd.PersistentFlags().String("region", "", "region")
	describeImagesCmd.PersistentFlags().String("profile", "", "profile")
	describeImagesCmd.PersistentFlags().String("repo", "", "repo")
}
