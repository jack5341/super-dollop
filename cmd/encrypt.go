package cmd

import (
	"fmt"
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
	gpgID = os.Getenv("GPG_ID")
	rootCmd.AddCommand(saveCmd)
	rootCmd.PersistentFlags().StringP("note", "n", "", "--note=here-is-my-note")
	rootCmd.PersistentFlags().StringP("file", "f", "", "--file=<YOUR FILE PATH>")
	rootCmd.PersistentFlags().BoolP("print", "p", false, "-p")
}

func encryptString(value string, isPrint bool) {
	cmd := exec.Command("gpg", "--encrypt", "-r", gpgID, "--armor")

	isDone, _ := pterm.DefaultSpinner.Start()

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

	if isPrint {
		fmt.Println(result)
	}

	defer isDone.Success("Successfully encrypted!")
}

func encryptFile(filePath string, isPrint bool) {
	cmd := exec.Command("gpg", "--encrypt", "--armor", "-r", gpgID, "-o", "/dev/stdout", filePath)

	// _ is result
	out, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	result := string(out)
	isDone, _ := pterm.DefaultSpinner.Start()

	if isPrint {
		fmt.Println(result)
	}

	/*
	if len(export) > 0 {
		err = ioutil.WriteFile( filePath + ".asc", []byte(result), 0644)
		if (err != nil) {
			log.Fatal(err)
		}
	}
	 */

	isDone.Success("Successfully encrypted!")
}