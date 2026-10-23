package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strings"
	"testing"

	internalprovider "github.com/CiscoDevNet/terraform-provider-mso/internal/provider"
	"github.com/CiscoDevNet/terraform-provider-mso/mso"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func TestMuxProviderServer(t *testing.T) {
	serverFactory, err := internalprovider.MuxProviderServer(context.Background(), "test", mso.Provider().GRPCProvider)
	if err != nil {
		t.Fatalf("creating mux provider server: %v", err)
	}

	response, err := serverFactory().GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("getting mux provider schema: %v", err)
	}

	for _, diagnostic := range response.Diagnostics {
		if diagnostic == nil || diagnostic.Severity != tfprotov6.DiagnosticSeverityError {
			continue
		}

		t.Fatalf("mux provider schema diagnostic: %s", diagnostic.Detail)
	}

	if _, ok := response.ResourceSchemas["mso_tenant"]; !ok {
		t.Fatal("mux provider did not expose existing SDKv2 resource mso_tenant")
	}
	if _, ok := response.DataSourceSchemas["mso_tenant"]; !ok {
		t.Fatal("mux provider did not expose existing SDKv2 data source mso_tenant")
	}
}

func TestLegacyLoggingPreservesSeverity(t *testing.T) {
	previousOutput, previousFlags := log.Writer(), log.Flags()
	t.Cleanup(func() {
		log.SetOutput(previousOutput)
		log.SetFlags(previousFlags)
	})

	var output bytes.Buffer
	initLogging(&output)
	log.SetFlags(log.LstdFlags)

	for _, level := range []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"} {
		t.Run(level, func(t *testing.T) {
			output.Reset()
			log.Printf("[%s] legacy API message", level)

			var entry struct {
				Level   string `json:"@level"`
				Message string `json:"@message"`
			}
			if err := json.Unmarshal(output.Bytes(), &entry); err != nil {
				t.Fatalf("expected structured log output: %v; got %q", err, output.String())
			}
			if entry.Level != strings.ToLower(level) {
				t.Errorf("expected severity %s, got %q", level, entry.Level)
			}
			if entry.Message != "legacy API message" {
				t.Errorf("expected timestamp and severity prefix to be removed, got %q", entry.Message)
			}
		})
	}
}
