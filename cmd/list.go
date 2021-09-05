package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your all encrypted files and notes.",
	Run: func(cmd *cobra.Command, args []string) {
		delete, _ := cmd.Flags().GetString("delete")
		if delete != "" {
			list(delete)
		}

		list("")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.PersistentFlags().StringP("delete", "d", "", "--delete=name-of-note-or-file")
}

func list(name string) {
	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	if name != "" {
		err := Client.RemoveObject(bucketName, name)
		if err != nil {
			panic(err)
		}
		fmt.Println("Removed Successfully!")
		return
	}

	// Notes
	NoteList := Client.ListObjects(bucketName, "notes", true, nil)

	var table [][]string

	NoteTable := append(table, []string{"Owner Display Name", "Key", "Size"})

	for v := range NoteList {
		size := strconv.FormatInt(v.Size, 10)
		NoteTable = append(NoteTable, []string{v.Owner.DisplayName, v.Key, size})
	}

	if len(NoteTable) > 1 {
		fmt.Print("Your Notes \n")
		fmt.Println("----------")
		pterm.DefaultTable.WithHasHeader().WithData(NoteTable).Render()
	}

	// Files
	FileList := Client.ListObjects(bucketName, "files", true, nil)

	FileTable := append(table, []string{"Owner Display Name", "Key", "Size"})

	for v := range FileList {
		size := strconv.FormatInt(v.Size, 10)
		FileTable = append(FileTable, []string{v.Owner.DisplayName, v.Key, size})
	}

	if len(FileList) > 1 {
		fmt.Print("Your Files \n")
		fmt.Println("----------")
		pterm.DefaultTable.WithHasHeader().WithData(FileTable).Render()
	}
}
