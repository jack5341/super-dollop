package cmd

import (
	"github.com/jack5341/super-dollop/pkg"
	"github.com/minio/minio-go"
	"github.com/spf13/cobra"
)

var cfgFile string
var Client *minio.Client

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dollop",
	Short: "A brief description of your application",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	Client = pkg.Connect()
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.super-dollop.yaml)")
}
