package cmd

import (
	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe a resource",
	Long: `Describe a resource
	For example:
	
	ecrctl describe repositories <repo-name>`,
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
