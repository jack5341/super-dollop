
package cmd

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"os/exec"
)

// noteCmd represents the note command
var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	godotenv.Load()
	rootCmd.AddCommand(noteCmd)
	getEncrypt("hello")
}

func getEncrypt(value string) string {
	gpgID := os.Getenv("GPG_ID")
	cmd := exec.Command("gpg" ,"--encrypt", "-r", gpgID, "--armor")

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

	return string(out)
}