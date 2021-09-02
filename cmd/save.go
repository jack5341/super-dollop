package cmd

import (
	"github.com/joho/godotenv"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"os/exec"
)

// saveCmd represents the note command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		note, _ := cmd.Flags().GetString("note")
		file, _ := cmd.Flags().GetString("file")
		if note != "" {
			encryptString(note)
		}

		if file != "" {
			encryptFile(file)
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
}

func encryptString(value string) {
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

	/*
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	result := string(out)
	 */
}

func encryptFile(filePath string) {
	cmd := exec.Command("gpg", "--encrypt", "--armor", "-r", gpgID, "-o", "/dev/stdout", filePath)

	isDone, _ := pterm.DefaultSpinner.Start("Encrypting...")

	defer isDone.Success("Successfully encrypted!")

	// _ is result
	_, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}
}