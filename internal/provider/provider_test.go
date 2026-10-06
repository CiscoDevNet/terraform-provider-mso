package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-mso/mso"
	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"mso": func() (tfprotov6.ProviderServer, error) {
		ctx := context.Background()

		upgradedSDKServer, err := tf5to6server.UpgradeServer(ctx, mso.Provider().GRPCProvider)
		if err != nil {
			return nil, err
		}

		providers := []func() tfprotov6.ProviderServer{
			providerserver.NewProtocol6(New("test")()),
			func() tfprotov6.ProviderServer {
				return upgradedSDKServer
			},
		}

		muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
		if err != nil {
			return nil, err
		}

		return muxServer.ProviderServer(), nil
	},
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
