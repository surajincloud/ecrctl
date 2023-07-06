package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource",
	Long: `Delete a resource
	For example:
	
	ecrctl delete repository <repo-name>`,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
