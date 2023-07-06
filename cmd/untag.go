/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var untagCmd = &cobra.Command{
	Use:   "untag",
	Short: "untag a resource",
	Long: `Untag a resource
	For example:
	
	ecrctl untag repository <repo-name>`,
}

func init() {
	rootCmd.AddCommand(untagCmd)
}
