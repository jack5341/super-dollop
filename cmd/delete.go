package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// listCmd represents the list command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete your files or notes.",
	Example: "dollop delete -n my-note\n# Or \ndollop delete -f my-file",
	Run: func(cmd *cobra.Command, args []string) {
		isNote, _ := cmd.Flags().GetBool("note")
		isFile, _ := cmd.Flags().GetBool("file")

		if args[0] == "" {
			log.Fatal("You need to provide file name or note name")
		}

		name := "/"
		if isNote {
			name += "note/"
		} else if isFile {
			name += "files/"
		}
		name += args[0] + ".asc"

		err := deleteObject(name)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolP("note", "n", true, "--note")
	deleteCmd.Flags().BoolP("file", "f", false, "--file")
}

func deleteObject(name string) error {
	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	err := Client.RemoveObject(bucketName, name)
	if err != nil {
		return errors.New("failed to remove object")
	}
	fmt.Println("removed successfully")

	return nil
}
