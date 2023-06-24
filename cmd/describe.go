package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe a resource",
	Long: `Describe a resource
	For example:
	
	ecrctl describe repositories <repo-name>`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("describe called")
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
