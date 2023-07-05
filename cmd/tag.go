/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "tag a resource",
	Long: `Describe a resource
	For example:
	
	ecrctl tag repository <repo-name>`,
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
