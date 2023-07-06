package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/spf13/cobra"
	awspkg "github.com/surajincloud/ecrctl/pkg/awsutil"
)

// deleteImagesCmd represents the image command
var deleteImageCmd = &cobra.Command{
	Use:     "images",
	Aliases: []string{"image"},
	Short:   "Delete an ECR image",
	Long: `Delete an ECR image
	For example:
		ecrctl delete image your-image
		`,
	RunE: deleteImage,
}

func deleteImage(cmd *cobra.Command, args []string) error {
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
		fmt.Println("please pass image name")
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

	_, err = client.BatchDeleteImage(ctx, &ecr.BatchDeleteImageInput{
		ImageIds: []types.ImageIdentifier{
			{
				ImageTag: &imageName,
			},
		},
	})

	if err != nil {
		return err
	}
	return nil
}
func init() {
	deleteCmd.AddCommand(deleteImageCmd)
	deleteCmd.PersistentFlags().String("region", "", "region")
	deleteCmd.PersistentFlags().String("profile", "", "profile")
	deleteCmd.PersistentFlags().String("image", "", "image")

}
