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
	aesEncrypt "crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/keybase/go-crypto/openpgp"
	"github.com/keybase/go-crypto/openpgp/armor"
	"github.com/spf13/cobra"
)

// encCmd represents the enc command
var encCmd = &cobra.Command{
	Use:   "enc",
	Short: "Encrypt your data AES256 or GPG key",
	Long: `Encrypt your data AES256 or GPG key.
Store them in minIO or S3 buckets.`,
	Run: func(cmd *cobra.Command, args []string) {
		note, _ := cmd.Flags().GetString("note")
		file, _ := cmd.Flags().GetString("file")
		save, _ := cmd.Flags().GetBool("save")
		gpg, _ := cmd.Flags().GetBool("gpg")
		aes, _ := cmd.Flags().GetBool("aes")
		password, _ := cmd.Flags().GetString("password")
		filename, _ := cmd.Flags().GetString("filename")

		// Environments
		gpgID := os.Getenv("DOLLOP_GPG_ID")
		endpoint := os.Getenv("MINIO_ENDPOINT")
		accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
		secretAccessKey := os.Getenv("MINIO_SECRET_KEY")

		if note != "" {
			if aes {
				if password == "" {
					fmt.Println("Please provide a password o encrypt your data with AES256")
					return
				}

				hasher := md5.New()
				hasher.Write([]byte(password))
				hasshedPass := hex.EncodeToString(hasher.Sum(nil))

				block, _ := aesEncrypt.NewCipher([]byte(hasshedPass))
				gcm, err := cipher.NewGCM(block)

				if err != nil {
					log.Fatal(err.Error())
					return
				}

				nonce := make([]byte, gcm.NonceSize())
				if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
					log.Fatal(err.Error())
					return
				}

				ciphertext := gcm.Seal(nonce, nonce, []byte(note), nil)

				if save {
					minioClient, err := minio.New(endpoint, &minio.Options{
						Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
						Secure: false,
					})

					if err != nil {
						log.Fatal(err.Error())
						return
					}

					if filename == "" {
						fmt.Println("If you want to save your encrypted data in minIO, please provide a filename")
						return
					}

					minioClient.PutObject(cmd.Context(), "dollop", filename, bytes.NewReader(ciphertext), int64(len(note)), minio.PutObjectOptions{})

					fmt.Println("Your data has been encrypted by AES256 and stored in minIO")
					return
				}

				fmt.Println(`Your note is encrypted by AES256: ` + string(ciphertext))
				return
			}

			if gpg {
				var err error
				var buffer *bytes.Buffer
				var armored io.WriteCloser
				var crypter io.WriteCloser

				entityList, _ := openpgp.ReadArmoredKeyRing(strings.NewReader(gpgID))

				buffer = bytes.NewBuffer(nil)

				armored, err = armor.Encode(buffer, note, nil)

				if err != nil {
					log.Fatal("Error while encoding note")
					return
				}

				crypter, _ = openpgp.Encrypt(armored, entityList, nil, nil, nil)

				crypter.Write([]byte(note))
				crypter.Close()

				if save {
					minioClient, err := minio.New(endpoint, &minio.Options{
						Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
						Secure: false,
					})

					if err != nil {
						log.Fatal(err.Error())
						return
					}

					if filename == "" {
						fmt.Println("If you want to save your encrypted data in minIO, please provide a filename")
						return
					}

					minioClient.PutObject(cmd.Context(), "dollop", filename, strings.NewReader(buffer.String()), int64(len(note)), minio.PutObjectOptions{})

					fmt.Println("Your data has been encrypted by GPG and stored in minIO")
					return
				}

				fmt.Println("Your note is encrypted by GPG: " + buffer.String())
				return
			}
		}

		if file != "" {
			var err error
			var buffer *bytes.Buffer
			var armored io.WriteCloser
			var crypter io.WriteCloser

			r := strings.NewReader(gpgID)

			entityList, _ := openpgp.ReadArmoredKeyRing(r)

			buffer = bytes.NewBuffer(nil)

			armored, err = armor.Encode(buffer, note, nil)

			if err != nil {
				log.Fatal("Error while encoding note")
			}

			crypter, _ = openpgp.Encrypt(armored, entityList, nil, nil, nil)

			crypter.Write([]byte(note))
			crypter.Close()

			fmt.Println(buffer.String())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(encCmd)
	encCmd.Flags().BoolP("gpg", "g", true, "Use GPG to encrypt your data")
	encCmd.Flags().BoolP("aes", "a", false, "Use AES256 to encrypt your data")
	encCmd.Flags().BoolP("save", "s", false, "Store your data in minio")
	encCmd.Flags().StringP("note", "n", "", "Encrypt plain text note")
	encCmd.Flags().StringP("file", "f", "", "Encrypt your file or directory")
	encCmd.Flags().StringP("password", "p", "", "Use your keyword as password to encrypt your input")
	encCmd.Flags().StringP("filename", "", "", "Filename to store your encrypted data")
}
