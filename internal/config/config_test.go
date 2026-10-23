package config

import "testing"

func TestBuildConfig(t *testing.T) {
	config, err := BuildConfig("user", "password", "https://example.test", false, "", "", "mso", "4")
	if err != nil {
		t.Fatalf("BuildConfig returned an error: %v", err)
	}

	if config.Username != "user" || config.Password != "password" || config.URL != "https://example.test" {
		t.Fatalf("unexpected credentials or URL: %#v", config)
	}
	if config.IsInsecure {
		t.Fatal("expected insecure=false to be preserved")
	}
	if config.MaxRetries != 4 {
		t.Fatalf("expected four retries, got %d", config.MaxRetries)
	}
}

func TestBuildConfigDefaultsRetries(t *testing.T) {
	config, err := BuildConfig("user", "password", "https://example.test", true, "", "", "", "")
	if err != nil {
		t.Fatalf("BuildConfig returned an error: %v", err)
	}
	if config.MaxRetries != 2 {
		t.Fatalf("expected two retries, got %d", config.MaxRetries)
	}
}
