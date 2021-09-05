package cmd

import (
	"github.com/spf13/cobra"
)

// decCommand represents the decrypt command
var decCommand = &cobra.Command{
	Use:   "list",
	Short: "List your all encrypted files and notes.",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(decCommand)
	rootCmd.Flags().StringP("name", "n", "", "The key to decrypt the file")
}
