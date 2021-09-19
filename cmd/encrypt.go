package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
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
		var err error
		if note != "" {
			err = encryptString(note, isPrint)
		}

		if file != "" {
			err = encryptFile(file, isPrint)
		}
		if err != nil {
			log.Fatal(err)
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

func encryptString(value string, isPrint bool) error {
	cmd := exec.Command("gpg", "--encrypt", "-r", gpgID, "--armor")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return errors.New("occurred with a problem while encrypt string")
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, value)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New("occurred with a problem while combined output")
	}

	result := string(out)

	if isPrint {
		fmt.Println(result)
		return nil
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
		return errors.New("occurred with a problem while upload encrypted string")
	}

	fmt.Println("successfully encrypted ", status)
	return nil
}

func encryptFile(filePath string, isPrint bool) error {
	cmd := exec.Command("gpg", "--encrypt", "--armor", "-r", gpgID, filePath)

	var err error
	var stdout io.Reader
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return errors.New("occurred with a problem while upload encrypted file")
	}
	defer cmd.Wait()

	if isPrint {
		stdout = io.TeeReader(stdout, os.Stdout)
	}

	fileName := filepath.Base(filePath)

	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	_, err = Client.PutObject(bucketName, "/files/"+fileName+".asc", stdout, -1, minio.PutObjectOptions{
		ContentType: "application/pgp-encrypted",
	})

	if err != nil {
		return errors.New("occurred with a problem while upload encrypted file")
	}

	if err := cmd.Wait(); err != nil {
		return errors.New("occurred with a problem while upload encrypted file")
	}

	fmt.Println("successfully encrypted")
	return nil
}
