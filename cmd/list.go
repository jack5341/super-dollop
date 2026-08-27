package cmd

import (
	"fmt"
	"strconv"

	"github.com/jack5341/super-dollop/internal/storage"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your encrypted notes and files",
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := getApp()
		if err != nil {
			return err
		}

		notes, err := app.client.List("notes/")
		if err != nil {
			return err
		}
		renderTable("Your Notes", notes)

		files, err := app.client.List("files/")
		if err != nil {
			return err
		}
		renderTable("Your Files", files)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func renderTable(title string, objects []storage.ObjectInfo) {
	if len(objects) == 0 {
		return
	}

	table := [][]string{{"Key", "Size"}}
	for _, obj := range objects {
		table = append(table, []string{obj.Key, strconv.FormatInt(obj.Size, 10)})
	}

	fmt.Printf("%s\n----------\n", title)
	pterm.DefaultTable.WithHasHeader().WithData(table).Render()
}
