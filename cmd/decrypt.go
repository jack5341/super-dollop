package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
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
		name, _ := cmd.Flags().GetString("name")
		var err error
		if name == "" {
			err = decrypt(name)
		}

		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(decCmd)
	decCmd.PersistentFlags().StringP("name", "n", "", "Decrypt and show your note or file")
}

func decrypt(name string) error {
	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	status, err := Client.GetObject(bucketName, name, minio.GetObjectOptions{})

	if err != nil {
		return errors.New("failed while getting data")
	}

	tempPath := "/tmp/" + path.Base(name)

	localFile, err := os.Create(tempPath)
	if err != nil {
		return errors.New("failed while creating file")
	}

	if _, err = io.Copy(localFile, status); err != nil {
		return errors.New("failed while copying data")
	}

	fmt.Println(tempPath)
	cmd := exec.Command("gpg", "--decrypt", tempPath)
	data, err := cmd.Output()
	if err != nil {
		return errors.New("failed while decrypting data")
	}

	defer fmt.Println(string(data))
	return nil
}
