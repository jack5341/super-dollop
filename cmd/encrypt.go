package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
	"github.com/minio/minio-go"
	"github.com/spf13/cobra"
)

// encCommand represents the encrypt command
var encCommand = &cobra.Command{
	Use:   "enc",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		note, _ := cmd.Flags().GetString("note")
		file, _ := cmd.Flags().GetString("file")
		isPrint, _ := cmd.Flags().GetBool("print")
		if note != "" {
			encryptString(note, isPrint)
		}

		if file != "" {
			encryptFile(file, isPrint)
		}
	},
}

var gpgID string

func init() {
	godotenv.Load()
	gpgID = os.Getenv("MINIO_GPG_ID")
	rootCmd.AddCommand(encCommand)
	encCommand.Flags().StringP("note", "n", "", "--note=here-is-my-note")
	encCommand.Flags().StringP("file", "f", "", "--file=<YOUR FILE PATH>")
	encCommand.Flags().BoolP("print", "p", false, "-p")
}

func encryptString(value string, isPrint bool) {
	cmd := exec.Command("gpg", "--encrypt", "-r", gpgID, "--armor")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err, "occurred with a problem while encrypt string")
		return
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, value)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err, "occurred with a problem while combined output")
		return
	}

	result := string(out)

	if isPrint {
		fmt.Println(result)
		return
	}

	validation := func(input string) error {
		if input == "" {
			return errors.New("note name is required input")
		}

		validPath := regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

		if !validPath.MatchString(input) {
			return errors.New("please enter valid file name")
		}

		if _, err := Client.StatObject(os.Getenv("MINIO_BUCKET_NAME"), "notes/"+input+".asc", minio.StatObjectOptions{}); err == nil {
			return errors.New("note is already exists")
		}

		return nil
	}

	notePrompt := promptui.Prompt{
		Label:    "Give a name to your note",
		Validate: validation,
	}

	noteName, _ := notePrompt.Run()

	readedResult := strings.NewReader(result)

	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	status, err := Client.PutObject(bucketName, "/notes/"+noteName+".asc", readedResult, int64(len(result)), minio.PutObjectOptions{
		ContentType: "application/pgp-encrypted",
	})

	if err != nil {
		fmt.Println(err, "occurred with a problem while upload encrypted string")
		return
	}

	fmt.Println("successfully encrypted ", status)
}

func encryptFile(filePath string, isPrint bool) {
	cmd := exec.Command("gpg", "--encrypt", "--armor", "-r", gpgID, "-o", "/dev/stdout", filePath)

	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err, "occurred with a problem while upload encrypted file")
		return
	}

	result := string(out)

	if isPrint {
		fmt.Println(result)
		return
	}

	readedResult := strings.NewReader(result)
	fileName := filepath.Base(filePath)

	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	status, err := Client.PutObject(bucketName, "/files/"+fileName+".asc", readedResult, int64(len(result)), minio.PutObjectOptions{
		ContentType: "application/pgp-encrypted",
	})

	if err != nil {
		fmt.Println(err, "occurred with a problem while upload encrypted file")
		return
	}

	fmt.Println("succesfully encrypted ", status)
}
