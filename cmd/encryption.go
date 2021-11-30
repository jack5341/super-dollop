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
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// encCmd represents the enc command
var encCmd = &cobra.Command{
	Use:   "enc",
	Short: "Encrypt your data SHA256 or GPG key",
	Long: `Encrypt your data SHA256 or GPG key.
Store them in minIO or S3 buckets.`,
	Run: func(cmd *cobra.Command, args []string) {
		note, _ := cmd.Flags().GetString("note")
		file, _ := cmd.Flags().GetString("file")
		gpgID := os.Getenv("DOLLOP_GPG_ID")

		// gpg, err := cmd.Flags().GetBool("gpg")
		// sha, err := cmd.Flags().GetBool("sha")
		// keep, err := cmd.Flags().GetBool("keep")

		if note != "" {
			cmd := exec.Command("gpg", "--encrypt", "-r", gpgID, "--armor")

			stdin, err := cmd.StdinPipe()

			if err != nil {
				errors.New("occurred with a problem while encrypt string")
			}

			go func() {
				defer stdin.Close()
				io.WriteString(stdin, note)
			}()

			out, err := cmd.CombinedOutput()
			if err != nil {
				errors.New("occurred with a problem while combined output")
			}

			result := string(out)

			fmt.Println("Encrypted note: ", result)
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
	encCmd.Flags().BoolP("gpg", "g", false, "Use GPG to encrypt your data")
	encCmd.Flags().BoolP("sha", "s", false, "Use SHA256 to encrypt your data")
	encCmd.Flags().BoolP("keep", "m", false, "Store your data in minio")
	encCmd.Flags().StringP("note", "n", "", "Encrypt plain text note")
	encCmd.Flags().StringP("file", "f", "", "Encrypt your file or directory")
}
