package storage

import (
	"fmt"
	"regexp"
)

// Kind selects the object prefix used for stored items.
type Kind int

const (
	Note Kind = iota
	File
)

var namePattern = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

func (k Kind) prefix() string {
	if k == File {
		return "files"
	}
	return "notes"
}

// Key builds the object key for a note or file name.
func Key(kind Kind, name string) string {
	return fmt.Sprintf("%s/%s.asc", kind.prefix(), name)
}

// ValidateName reports whether name is a valid note/file name.
func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if !namePattern.MatchString(name) {
		return fmt.Errorf("invalid name %q: only letters, numbers, underscores and dashes are allowed", name)
	}
	return nil
}
