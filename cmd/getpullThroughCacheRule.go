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

// pullThroughCacheRuleCmd represents the pullThroughCacheRule command
var getpullThroughCacheRuleCmd = &cobra.Command{
	Use:     "pullThroughCacheRule",
	Short:   "Get ECR Pull through cache rules",
	Aliases: []string{"ptc"},
	Long: `Get ECR Pull through cache rules
    Examples:

	# List all Pull through cache rules
	ecrctl get pullThroughCacheRule

	# List all repositories their tags
	ecrctl get pullThroughCacheRule <your-pull-through-cache-rule>
	`,
	RunE: getpullThroughCacheRule,
}

func getpullThroughCacheRule(cmd *cobra.Command, args []string) error {

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

	// aws config
	cfg, err := awspkg.GetAWSConfig(ctx, region, profile)
	if err != nil {
		log.Fatal(err)
	}
	// Create an EKS client using the loaded configuration
	client := ecr.NewFromConfig(cfg)

	ptcrs, err := client.DescribePullThroughCacheRules(ctx, &ecr.DescribePullThroughCacheRulesInput{})
	if err != nil {
		log.Fatal(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 5, 2, 3, ' ', tabwriter.TabIndent)
	defer w.Flush()
	fmt.Fprintln(w, "NAME", "\t", "URL", "\t", "AGE")

	for _, i := range ptcrs.PullThroughCacheRules {
		fmt.Fprintln(w, *i.EcrRepositoryPrefix, "\t", *i.UpstreamRegistryUrl, "\t", awspkg.GetAge(*i.CreatedAt))
	}

	return nil
}
func init() {
	getCmd.AddCommand(getpullThroughCacheRuleCmd)
	getpullThroughCacheRuleCmd.PersistentFlags().String("region", "", "region")
	getpullThroughCacheRuleCmd.PersistentFlags().String("profile", "", "profile")

}
