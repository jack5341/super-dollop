package cmd

import (
	"fmt"
	"os"

	"github.com/minio/minio-go"
	"github.com/spf13/cobra"
)

// decCmd represents the decrypt command
var decCmd = &cobra.Command{
	Use:   "dec",
	Short: "List your all encrypted files and notes.",
	Run: func(cmd *cobra.Command, args []string) {
		decrypt(cmd.Flag("fname").Value.String(), cmd.Flag("name").Value.String())
	},
}

func init() {
	rootCmd.AddCommand(decCmd)
	decCmd.Flags().StringP("fname", "f", "", "Decrypt and show your file")
	decCmd.Flags().StringP("name", "n", "", "Decrypt and show your note")
}

func decrypt(filename string, notename string) {
	fmt.Println(filename, notename)
	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	status, err := Client.GetObject(bucketName, "notes/note-1.asc", minio.GetObjectOptions{})

	if err != nil {
		panic(err)
	}

	fmt.Print(status)
}
