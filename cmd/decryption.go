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
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/keybase/go-crypto/openpgp"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
)

// decryptionCmd represents the decryption command
var decryptionCmd = &cobra.Command{
	Use:   "dec",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		// Flags
		store, _ := cmd.Flags().GetBool("store")
		filename, _ := cmd.Flags().GetString("filename")

		// Environments
		gpgID := os.Getenv("DOLLOP_GPG_ID")
		endpoint := os.Getenv("MINIO_ENDPOINT")
		accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
		secretAccessKey := os.Getenv("MINIO_SECRET_KEY")

		if store {
			minioClient, err := minio.New(endpoint, &minio.Options{
				Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
				Secure: false,
			})

			if err != nil {
				log.Fatal(err.Error())
				return
			}

			if filename == "" {
				fmt.Println("if you want to decrypt your data, please provide a filename or path")
				return
			}

			reader, _ := minioClient.GetObject(cmd.Context(), "dollop", filename, minio.GetObjectOptions{})

			data, err := ioutil.ReadAll(reader)

			if err != nil {
				log.Fatal(err.Error())
				return
			}

			var entityList openpgp.EntityList

			keyringFileBuffer, err := os.Open(gpgID)
			if err != nil {
				log.Fatal(err.Error())
				return
			}

			defer keyringFileBuffer.Close()

			// Fix: open GPG ID: no such file or directory
			// While decryiption.

			entityList, err = openpgp.ReadKeyRing(keyringFileBuffer)
			if err != nil {
				log.Fatal(err.Error())
				return
			}

			md, err := openpgp.ReadMessage(bytes.NewBuffer(data), entityList, nil, nil)

			if err != nil {
				log.Fatal(err.Error())
				return
			}

			bytes, err := ioutil.ReadAll(md.UnverifiedBody)

			if err != nil {
				log.Fatal(err.Error())
				return
			}

			decStr := string(bytes)

			fmt.Println(decStr)
		}
	},
}

func init() {
	rootCmd.AddCommand(decryptionCmd)
	decryptionCmd.Flags().Bool("store", false, "define if the name of file in the store.")
	decryptionCmd.Flags().String("filename", "", "name of file for want to decrypt")
}
