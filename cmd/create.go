package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource",
	Long: `Describe a resource
	For example:
	
	ecrctl create repository <repo-name>`,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
