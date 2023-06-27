package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/spf13/cobra"
	awspkg "github.com/surajincloud/ecrctl/pkg/awsutil"
)

// repositoriesCmd represents the repositories command
var getRepositoriesCmd = &cobra.Command{
	Use:     "repositories",
	Short:   "Get ECR repositories",
	Aliases: []string{"repo"},
	Long: `List ECR repositories
    Examples:

	# List all repositories with URI
	ecrctl get repositories

	# List all repositories their tags
	ecrctl get repositories --show-tags

	# List all repositories with given tags
	ecrctl get repositories --tag key=value

	or 

	ecrctl get repositories -t key=value
	`,
	RunE: repositories,
}

func repositories(cmd *cobra.Command, args []string) error {

	ctx := context.Background()
	// read flag values
	region, err := cmd.Flags().GetString("region")
	if err != nil {
		fmt.Println("no region", err)
	}

	tags, err := cmd.Flags().GetString("tag")
	if err != nil {
		fmt.Println("no tags", err)
	}

	showTags, err := cmd.Flags().GetBool("show-tags")
	if err != nil {
		fmt.Println("no tags", err)
	}

	// aws config
	cfg, err := awspkg.GetAWSConfig(ctx, region)
	if err != nil {
		log.Fatal(err)
	}
	// Create an EKS client using the loaded configuration
	client := ecr.NewFromConfig(cfg)

	repoList, err := awspkg.ListRepos(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 5, 2, 3, ' ', tabwriter.TabIndent)
	defer w.Flush()

	if !showTags && tags == "" {
		fmt.Fprintln(w, "NAME", "\t", "URI")
		for _, i := range repoList {
			fmt.Fprintln(w, i.Name, "\t", i.Uri)
		}
	}

	if !showTags && tags != "" {
		fmt.Fprintln(w, "NAME", "\t", "URI")
		for _, i := range repoList {
			var tagList []string
			for _, k := range i.Tags {
				tagList = append(tagList, fmt.Sprintf("%s=%s", k.Key, k.Value))
			}
			tagsStringFormat := strings.Join(tagList, ",")
			if tags != "" {
				if strings.Contains(tagsStringFormat, tags) {
					fmt.Fprintln(w, i.Name, "\t", i.Uri)
				}
			} else {
				fmt.Fprintln(w, i.Name, "\t", i.Uri)
			}
		}
	}

	if showTags && tags == "" {
		fmt.Fprintln(w, "NAME", "\t", "URI", "\t", "TAGS")
		for _, i := range repoList {
			var tagList []string
			for _, k := range i.Tags {
				tagList = append(tagList, fmt.Sprintf("%s=%s", k.Key, k.Value))
			}
			tagsStringFormat := strings.Join(tagList, ",")
			if tags != "" {
				if strings.Contains(tagsStringFormat, tags) {
					fmt.Fprintln(w, i.Name, "\t", i.Uri, "\t", tagsStringFormat)
				}
			} else {
				fmt.Fprintln(w, i.Name, "\t", i.Uri, "\t", tagsStringFormat)
			}
		}
	}

	return err
}

func init() {
	getCmd.AddCommand(getRepositoriesCmd)
	getCmd.PersistentFlags().String("region", "", "region")
	getCmd.PersistentFlags().StringP("tag", "t", "", "tags")
	getCmd.PersistentFlags().Bool("show-tags", false, "show tags")
}
