/*
Copyright © 2023 Derek Worthen <worthend.derek@gmail.com>
*/
package cmd

import (
	_ "embed"

	"github.com/dworthen/goscripty/generators"
	"github.com/spf13/cobra"
)

// githubCmd represents the github command
var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "Generate install scripts for github artifacts",
	Long:  `Generate install scripts for github artifacts`,
	Run: func(cmd *cobra.Command, args []string) {
		generators.GenerateGithubInstallers()
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// githubCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// githubCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
