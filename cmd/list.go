/*
Copyright © 2021 NEDIM AKAR nedim.akar53411@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"log"
	"os"

	util "github.com/jack5341/super-dollop/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type TableStruct struct {
	Name         string
	Size         int64
	LastModified string
	ContentType  string
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		detail, _ := cmd.Flags().GetBool("detail")

		// Environments
		endpoint := os.Getenv("MINIO_ENDPOINT")
		accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
		secretAccessKey := os.Getenv("MINIO_SECRET_KEY")

		if endpoint == "" || accessKeyID == "" || secretAccessKey == "" {
			log.Fatal("Please set MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY")
			return
		}

		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: false,
		})

		if err != nil {
			log.Fatal(err)
			return
		}

		list := minioClient.ListObjects(cmd.Context(), "dollop", minio.ListObjectsOptions{})

		if detail {
			var table [][]string
			listTable := append(table, []string{"Element Name", "Last Modified", "Size"})

			for e := range list {
				sizeOfFile := util.ConvertByte(float64(e.Size))
				listTable = append(listTable, []string{e.Key, e.LastModified.String(), sizeOfFile})
			}

			pterm.DefaultTable.WithHasHeader().WithData(listTable).Render()
			return
		}

		for e := range list {
			sizeOfFile := util.ConvertByte(float64(e.Size))
			nameOfFile := e.Key
			fmt.Println(nameOfFile + " ----- " + sizeOfFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().Bool("detail", false, "List encrypted files with details")
}
