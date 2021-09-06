package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/minio/minio-go"
	"github.com/spf13/cobra"
)

// decCmd represents the decrypt command
var decCmd = &cobra.Command{
	Use:   "dec",
	Short: "List your all encrypted files and notes.",
	Run: func(cmd *cobra.Command, args []string) {
		decrypt(cmd.Flag("name").Value.String())
	},
}

func init() {
	rootCmd.AddCommand(decCmd)
	decCmd.Flags().StringP("name", "n", "", "Decrypt and show your note or file")
}

func decrypt(name string) {
	if name == "" {
		panic("Please provide a name")
	}

	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	status, err := Client.GetObject(bucketName, name, minio.GetObjectOptions{})

	if err != nil {
		panic(err)
	}

	tempPath := "/tmp/" + path.Base(name)

	localFile, err := os.Create(tempPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err = io.Copy(localFile, status); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(tempPath)
	cmd := exec.Command("gpg", "--decrypt", tempPath)
	data, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
