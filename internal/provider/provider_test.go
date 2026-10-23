package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMSOProviderModelConfigUsesEnvironment(t *testing.T) {
	t.Setenv("MSO_USERNAME", "user")
	t.Setenv("MSO_PASSWORD", "password")
	t.Setenv("MSO_URL", "https://example.test")
	t.Setenv("MSO_INSECURE", "false")
	t.Setenv("MSO_RETRIES", "4")

	config, err := (MSOProviderModel{
		Username: types.StringNull(),
		Password: types.StringNull(),
		URL:      types.StringNull(),
		Insecure: types.BoolNull(),
		Retries:  types.StringNull(),
	}).Config()
	if err != nil {
		t.Fatalf("configuring Framework provider: %v", err)
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
