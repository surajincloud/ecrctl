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
	Use:   "images",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: describeImages,
}

func describeImages(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	// read flag values
	region, err := cmd.Flags().GetString("region")
	if err != nil {
		fmt.Println("no region", err)
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
	cfg, err := awspkg.GetAWSConfig(ctx, region)
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
	// fmt.Println("enhanced")
	// for _, i := range image.EnhancedScanFindings {
	// 	fmt.Println("test")
	// 	fmt.Fprintln(w, *i.Status, "\t", i.Severity, "\t", *i.ScoreDetails)

	// }
	return nil
}

func init() {
	describeCmd.AddCommand(describeImagesCmd)
	describeImagesCmd.PersistentFlags().String("region", "", "region")
	describeImagesCmd.PersistentFlags().String("repo", "", "repo")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
