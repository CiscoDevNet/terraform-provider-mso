package main

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func TestMuxProviderServer(t *testing.T) {
	serverFactory, err := muxProviderServer(context.Background())
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
}
