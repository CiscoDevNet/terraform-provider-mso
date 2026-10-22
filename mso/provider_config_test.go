package mso

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestConfigFromSDKResourceDataUsesEnvironment(t *testing.T) {
	t.Setenv("MSO_USERNAME", "user")
	t.Setenv("MSO_PASSWORD", "password")
	t.Setenv("MSO_URL", "https://example.test")
	t.Setenv("MSO_INSECURE", "false")
	t.Setenv("MSO_RETRIES", "4")

	data := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{})
	config, err := configFromSDKResourceData(data)
	if err != nil {
		t.Fatalf("configuring SDKv2 provider: %v", err)
	}

	if config.Username != "user" || config.Password != "password" || config.URL != "https://example.test" {
		t.Fatalf("unexpected credentials or URL: %#v", config)
	}
	if config.IsInsecure {
		t.Fatal("expected MSO_INSECURE=false to be respected")
	}
	if config.MaxRetries != 4 {
		t.Fatalf("expected four retries, got %d", config.MaxRetries)
	}
}

func TestProviderConfigDefaultsInsecureAndRetries(t *testing.T) {
	t.Setenv("MSO_USERNAME", "user")
	t.Setenv("MSO_PASSWORD", "password")
	t.Setenv("MSO_URL", "https://example.test")
	t.Setenv("MSO_INSECURE", "")
	t.Setenv("MSO_RETRIES", "")

	data := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{})
	config, err := configFromSDKResourceData(data)
	if err != nil {
		t.Fatalf("configuring SDKv2 provider: %v", err)
	}

	if !config.IsInsecure {
		t.Fatal("expected insecure client behavior by default")
	}
	if config.MaxRetries != 2 {
		t.Fatalf("expected two retries by default, got %d", config.MaxRetries)
	}
}
