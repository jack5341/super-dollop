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
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

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
		gpg, _ := cmd.Flags().GetBool("gpg")
		aes, _ := cmd.Flags().GetBool("aes")
		password, _ := cmd.Flags().GetString("password")

		gpgID := os.Getenv("DOLLOP_GPG_ID")

		if note != "" {
			var err error
			var buffer *bytes.Buffer
			var armored io.WriteCloser
			var crypter io.WriteCloser

			if aes {
				if password == "" {
					fmt.Println("Please provide a password to encrypt your data with AES256")
					return
				}

				hasher := md5.New()
				hasher.Write([]byte(password))
				hasshedPass := hex.EncodeToString(hasher.Sum(nil))

				block, _ := aesEncrypt.NewCipher([]byte(hasshedPass))
				gcm, err := cipher.NewGCM(block)

				if err != nil {
					errors.New(err.Error())
				}

				nonce := make([]byte, gcm.NonceSize())
				if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
					errors.New(err.Error())
				}

				ciphertext := gcm.Seal(nonce, nonce, []byte(note), nil)
				fmt.Println(ciphertext)
				return
			}

			if gpg {
				r := strings.NewReader(gpgID)

				entityList, _ := openpgp.ReadArmoredKeyRing(r)

				buffer = bytes.NewBuffer(nil)

				armored, err = armor.Encode(buffer, note, nil)

				if err != nil {
					errors.New("Error while encoding note")
				}

				crypter, _ = openpgp.Encrypt(armored, entityList, nil, nil, nil)

				crypter.Write([]byte(note))
				crypter.Close()

				fmt.Println(buffer.String())
				return
			}

		}

		if file != "" {
			cmd := exec.Command("gpg", "--encrypt", "--armor", "-r", gpgID, file)

			var err error
			var stdout io.Reader
			if stdout, err = cmd.StdoutPipe(); err != nil {
				errors.New("occurred with a problem while encrypting file")
			}

			err = cmd.Start()
			if err != nil {
				errors.New("occurred with a problem while upload encrypted file")
			}

			defer cmd.Wait()

			fmt.Println("successfully encrypted", stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(encCmd)
	encCmd.Flags().BoolP("gpg", "g", true, "Use GPG to encrypt your data")
	encCmd.Flags().BoolP("aes", "a", false, "Use AES256 to encrypt your data")
	encCmd.Flags().BoolP("keep", "m", false, "Store your data in minio")
	encCmd.Flags().StringP("note", "n", "", "Encrypt plain text note")
	encCmd.Flags().StringP("file", "f", "", "Encrypt your file or directory")
	encCmd.Flags().StringP("password", "p", "", "Use your keyword as password to encrypt your input")
}
