package config

import (
	"strings"
	"testing"
)

func TestValidateMissingVars(t *testing.T) {
	cfg := &Config{}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("Validate() expected error for empty config, got nil")
	}

	for _, want := range []string{"MINIO_ENDPOINT", "MINIO_ACCESS_KEY", "MINIO_SECRET_KEY", "MINIO_BUCKET_NAME", "MINIO_GPG_ID"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error should mention %s, got: %v", want, err)
		}
	}
}

func TestValidateComplete(t *testing.T) {
	cfg := &Config{
		Endpoint:  "127.0.0.1:9000",
		AccessKey: "admin",
		SecretKey: "secretadmin",
		Bucket:    "dollop-files",
		GPGID:     "5355D0DCF50CD0A8F3DE238BE1E313D9BB2B0ACE",
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
	}
}

func TestValidateSingleMissing(t *testing.T) {
	cfg := &Config{
		Endpoint:  "127.0.0.1:9000",
		AccessKey: "admin",
		SecretKey: "secretadmin",
		Bucket:    "dollop-files",
	}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("Validate() expected error for missing GPGID, got nil")
	}
	if !strings.Contains(err.Error(), "MINIO_GPG_ID") {
		t.Errorf("error should mention MINIO_GPG_ID, got: %v", err)
	}
}
