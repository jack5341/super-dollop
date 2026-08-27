package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/jack5341/super-dollop/internal/config"
	"github.com/jack5341/super-dollop/internal/storage"
	"github.com/spf13/cobra"
)

// version is set at build time via -ldflags "-X .../cmd.version=...".
var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "dollop",
	Short: "Encrypt notes and files with GPG and store them in MinIO/S3",
	Long: `dollop encrypts your notes and files with your own GPG key and
stores them in a MinIO or S3 object storage, keeping them safe and
portable. You can also print an encrypted note or file straight to
the terminal.`,
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command. It is called by main.main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

type appState struct {
	cfg    *config.Config
	client *storage.Client
}

var (
	appOnce sync.Once
	app     *appState
	appErr  error
)

// getApp lazily loads config and connects to MinIO, so commands that
// do not need storage (help, --version) work offline.
func getApp() (*appState, error) {
	appOnce.Do(func() {
		cfg, err := config.Load()
		if err != nil {
			appErr = err
			return
		}
		client, err := storage.NewClient(cfg)
		if err != nil {
			appErr = err
			return
		}
		app = &appState{cfg: cfg, client: client}
	})
	return app, appErr
}
