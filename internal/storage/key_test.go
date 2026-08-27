package storage

import "testing"

func TestKey(t *testing.T) {
	tests := []struct {
		kind Kind
		name string
		want string
	}{
		{Note, "my-note", "notes/my-note.asc"},
		{Note, "a", "notes/a.asc"},
		{File, "report.pdf", "files/report.pdf.asc"},
	}
	for _, tt := range tests {
		if got := Key(tt.kind, tt.name); got != tt.want {
			t.Errorf("Key(%v, %q) = %q, want %q", tt.kind, tt.name, got, tt.want)
		}
	}
}

func TestValidateName(t *testing.T) {
	valid := []string{"note", "my-note", "my_note", "note-123", "A"}
	for _, name := range valid {
		if err := ValidateName(name); err != nil {
			t.Errorf("ValidateName(%q) unexpected error: %v", name, err)
		}
	}

	invalid := []string{"", "has space", "dot.name", "slash/name", "unicode-ü", "tab\tname"}
	for _, name := range invalid {
		if err := ValidateName(name); err == nil {
			t.Errorf("ValidateName(%q) expected error, got nil", name)
		}
	}
}
