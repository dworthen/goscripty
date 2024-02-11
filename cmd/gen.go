package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:     "gen",
	Aliases: []string{"generate"},
	Short:   "Generate install scripts",
	Long:    `Generate install scripts`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := rootCmd.Help(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(1)
	},
}

func init() {
	genCmd.AddCommand(githubCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
