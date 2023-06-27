package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	Long: `get command
    For example:
    
	ecrctl get repositories
	`,
}

func init() {
	rootCmd.AddCommand(getCmd)
}
