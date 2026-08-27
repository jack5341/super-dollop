package cmd

import (
	"fmt"

	"github.com/jack5341/super-dollop/internal/storage"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a stored note or file",
	Example: `  dollop delete -n my-note
  dollop delete -f my-file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		note, _ := cmd.Flags().GetString("note")
		file, _ := cmd.Flags().GetString("file")

		if (note == "") == (file == "") {
			return fmt.Errorf("provide exactly one of -n/--note or -f/--file")
		}

		app, err := getApp()
		if err != nil {
			return err
		}

		kind := storage.Note
		name := note
		if file != "" {
			kind = storage.File
			name = file
		}

		key := storage.Key(kind, name)
		if err := app.client.Remove(key); err != nil {
			return err
		}

		fmt.Printf("removed %s\n", key)
		return nil
	},
}

func init() {
	deleteCmd.Flags().StringP("note", "n", "", "name of the note to delete")
	deleteCmd.Flags().StringP("file", "f", "", "name of the file to delete")
	rootCmd.AddCommand(deleteCmd)
}
