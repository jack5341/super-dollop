package gpg

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// EncryptString encrypts value for recipient and returns the armored
// PGP output.
func EncryptString(recipient, value string) (string, error) {
	cmd := exec.Command("gpg", "--batch", "--yes", "--quiet", "--armor", "--encrypt", "-r", recipient)
	cmd.Stdin = strings.NewReader(value)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("gpg encrypt: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return stdout.String(), nil
}

// EncryptFile starts gpg encryption of path and returns a stream of
// the armored output plus a finish function that waits for gpg to
// exit. The caller must consume the stream (e.g. upload it) before
// calling finish.
func EncryptFile(recipient, path string) (io.ReadCloser, func() error, error) {
	cmd := exec.Command("gpg", "--batch", "--yes", "--quiet", "--armor", "--encrypt", "-r", recipient, path)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("gpg pipe: %w", err)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("gpg start: %w", err)
	}

	finish := func() error {
		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("gpg encrypt: %w: %s", err, strings.TrimSpace(stderr.String()))
		}
		return nil
	}

	return stdout, finish, nil
}

// Decrypt decrypts the armored file at path and returns the plaintext.
func Decrypt(path string) ([]byte, error) {
	cmd := exec.Command("gpg", "--batch", "--yes", "--quiet", "--decrypt", path)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("gpg decrypt: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return stdout.Bytes(), nil
}
