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
	"github.com/spf13/cobra"
)

// encCmd represents the enc command
var encCmd = &cobra.Command{
	Use:   "enc",
	Short: "Encrypt your data SHA256 or GPG key",
	Long: `Encrypt your data SHA256 or GPG key.
Store them in minio or S3 buckets.`,
	Run: func(cmd *cobra.Command, args []string) {
		gpg, _ := cmd.Flags().GetBool("gpg")
		sha, _ := cmd.Flags().GetBool("sha")
		keep, _ := cmd.Flags().GetBool("keep")

		var err error

		if gpg {

		}

		if sha {

		}
	},
}

func init() {
	rootCmd.AddCommand(encCmd)
	encCmd.Flags().BoolP("gpg", "g", false, "Use GPG to encrypt your data")
	encCmd.Flags().BoolP("sha", "s", false, "Use SHA256 to encrypt your data")
	encCmd.Flags().BoolP("keep", "m", false, "Store your data in minio")
}

func encryptWithGPG() {

}

func encryptWithSHA() {

}

func encryptWithKeep() {

}
