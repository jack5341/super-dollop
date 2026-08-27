package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/jack5341/super-dollop/internal/gpg"
	"github.com/jack5341/super-dollop/internal/storage"
	"github.com/spf13/cobra"
)

var decCmd = &cobra.Command{
	Use:   "dec",
	Short: "Decrypt a stored note or file and print it",
	Example: `  dollop dec -n my-note
  dollop dec -f my-file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		note, _ := cmd.Flags().GetString("note")
		file, _ := cmd.Flags().GetString("file")

		if (note == "") == (file == "") {
			return fmt.Errorf("provide exactly one of -n/--note or -f/--file")
		}

		kind := storage.Note
		name := note
		if file != "" {
			kind = storage.File
			name = file
		}

		return decrypt(kind, name)
	},
}

func init() {
	decCmd.Flags().StringP("note", "n", "", "name of the note to decrypt")
	decCmd.Flags().StringP("file", "f", "", "name of the file to decrypt")
	rootCmd.AddCommand(decCmd)
}

func decrypt(kind storage.Kind, name string) error {
	app, err := getApp()
	if err != nil {
		return err
	}

	key := storage.Key(kind, name)
	obj, err := app.client.Get(key)
	if err != nil {
		return err
	}
	defer obj.Close()

	tmp, err := os.CreateTemp("", "dollop-*.asc")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if _, err := io.Copy(tmp, obj); err != nil {
		tmp.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	plaintext, err := gpg.Decrypt(tmpPath)
	if err != nil {
		return err
	}

	if _, err := os.Stdout.Write(plaintext); err != nil {
		return fmt.Errorf("write plaintext: %w", err)
	}
	if len(plaintext) > 0 && plaintext[len(plaintext)-1] != '\n' {
		fmt.Println()
	}

	return nil
}
