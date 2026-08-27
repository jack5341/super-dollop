package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jack5341/super-dollop/internal/gpg"
	"github.com/jack5341/super-dollop/internal/storage"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var encCmd = &cobra.Command{
	Use:   "enc",
	Short: "Encrypt a note or file and store it",
	Example: `  dollop enc -n "here-is-my-note"
  dollop enc -f ./secret.txt
  dollop enc -n "quick note" -p`,
	RunE: func(cmd *cobra.Command, args []string) error {
		note, _ := cmd.Flags().GetString("note")
		file, _ := cmd.Flags().GetString("file")
		isPrint, _ := cmd.Flags().GetBool("print")

		if (note == "") == (file == "") {
			return fmt.Errorf("provide exactly one of -n/--note or -f/--file")
		}

		if note != "" {
			return encryptNote(note, isPrint)
		}
		return encryptFile(file, isPrint)
	},
}

func init() {
	encCmd.Flags().StringP("note", "n", "", "text to encrypt and store as a note")
	encCmd.Flags().StringP("file", "f", "", "path to a file to encrypt and store")
	encCmd.Flags().BoolP("print", "p", false, "print the encrypted output to the terminal")
	rootCmd.AddCommand(encCmd)
}

func encryptNote(note string, isPrint bool) error {
	app, err := getApp()
	if err != nil {
		return err
	}

	armored, err := gpg.EncryptString(app.cfg.GPGID, note)
	if err != nil {
		return err
	}

	if isPrint {
		fmt.Println(armored)
		return nil
	}

	name, err := promptForNoteName(app.client)
	if err != nil {
		return err
	}

	key := storage.Key(storage.Note, name)
	if err := app.client.Put(key, strings.NewReader(armored), int64(len(armored))); err != nil {
		return err
	}

	fmt.Printf("note %q encrypted and stored at %s\n", name, key)
	return nil
}

func promptForNoteName(client *storage.Client) (string, error) {
	prompt := promptui.Prompt{
		Label: "Give a name to your note",
		Validate: func(input string) error {
			if err := storage.ValidateName(input); err != nil {
				return err
			}
			if client.Stat(storage.Key(storage.Note, input)) == nil {
				return fmt.Errorf("note %q already exists", input)
			}
			return nil
		},
	}

	name, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return "", fmt.Errorf("aborted")
		}
		return "", fmt.Errorf("prompt: %w", err)
	}

	return name, nil
}

func encryptFile(path string, isPrint bool) error {
	app, err := getApp()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	stdout, finish, err := gpg.EncryptFile(app.cfg.GPGID, path)
	if err != nil {
		return err
	}
	defer stdout.Close()

	reader := io.Reader(stdout)
	if isPrint {
		reader = io.TeeReader(stdout, os.Stdout)
	}

	key := storage.Key(storage.File, filepath.Base(path))
	if err := app.client.Put(key, reader, -1); err != nil {
		return err
	}

	if err := finish(); err != nil {
		return err
	}

	fmt.Printf("file %q encrypted and stored at %s\n", filepath.Base(path), key)
	return nil
}
