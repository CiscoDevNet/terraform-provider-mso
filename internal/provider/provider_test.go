package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-mso/mso"
	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func testAccProtoV6ProviderFactories(t *testing.T) map[string]func() (tfprotov6.ProviderServer, error) {
	t.Helper()
	return map[string]func() (tfprotov6.ProviderServer, error){
		"mso": func() (tfprotov6.ProviderServer, error) {
			factory, err := MuxProviderServer(context.Background(), "test", mso.Provider().GRPCProvider)
			if err != nil {
				return nil, err
			}
			server := factory()
			// The test harness closes v6 transports without stopping the SDK server.
			t.Cleanup(func() {
				response, err := server.StopProvider(context.Background(), &tfprotov6.StopProviderRequest{})
				if err != nil {
					t.Errorf("stop mux provider: %v", err)
					return
				}
				if response != nil && response.Error != "" {
					t.Errorf("stop mux provider: %s", response.Error)
				}
			})
			return server, nil
		},
	}
}

func testAccProviderPreCheck(t *testing.T) {
	t.Helper()
	for _, name := range []string{"MSO_USERNAME", "MSO_PASSWORD", "MSO_URL"} {
		if os.Getenv(name) == "" {
			t.Fatalf("%s must be set for MSO acceptance tests", name)
		}
	}
}

func testAccSiteName() string {
	if siteName := os.Getenv("MSO_SITE_NAME1"); siteName != "" {
		return siteName
	}
	return "ansible_test"
}

func testAccResourceName(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
}

func testAccAPIClient() *client.Client {
	platform := os.Getenv("MSO_PLATFORM")
	if platform == "" {
		platform = "mso"
	}

	maxRetries := 2
	if retries := os.Getenv("MSO_RETRIES"); retries != "" {
		if parsedRetries, err := strconv.Atoi(retries); err == nil {
			maxRetries = parsedRetries
		}
	}

	return client.GetClient(
		os.Getenv("MSO_URL"),
		os.Getenv("MSO_USERNAME"),
		client.Password(os.Getenv("MSO_PASSWORD")),
		client.Insecure(true),
		client.Platform(platform),
		client.MaxRetries(maxRetries),
	)
}

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
