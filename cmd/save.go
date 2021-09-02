package cmd

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// saveCmd represents the note command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		note, _ := cmd.Flags().GetString("note")
		file, _ := cmd.Flags().GetString("file")
		export, _ := cmd.Flags().GetString("export")
		if note != "" {
			encryptString(note,export)
		}

		if file != "" {
			encryptFile(file,export)
		}
	},
}

var gpgID string

func init() {
	godotenv.Load()
	gpgID = os.Getenv("GPG_ID")
	rootCmd.AddCommand(saveCmd)
	rootCmd.PersistentFlags().StringP("note", "n", "", "--note=here-is-my-note")
	rootCmd.PersistentFlags().StringP("file", "f", "", "--file=<YOUR FILE PATH>")
	rootCmd.PersistentFlags().StringP("export", "e", ".", "--export=<EXPORT PATH>")
}

func encryptString(value string, edit string) {
	cmd := exec.Command("gpg", "--encrypt", "-r", gpgID, "--armor")

	isDone, _ := pterm.DefaultSpinner.Start("Encrypting...")

	defer isDone.Success("Successfully encrypted!")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, value)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	result := string(out)
	fmt.Println(result)
}

func encryptFile(filePath string, export string) {
	cmd := exec.Command("gpg", "--encrypt", "--armor", "-r", gpgID, "-o", "/dev/stdout", filePath)

	// _ is result
	out, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	result := string(out)
	isDone, _ := pterm.DefaultSpinner.Start("Encrypting...")

	path := filePath
	fileName := filepath.Base(path)

	if len(export) > 0 {
		err = ioutil.WriteFile(fileName + ".asc", []byte(result), 0644)
		if(err != nil) {
			log.Fatal(err)
		}
	}

	defer isDone.Success("Successfully encrypted!")
}