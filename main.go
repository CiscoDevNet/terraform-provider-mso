package main

import (
	"context"
	"io"
	"log"
	"os"

	internalprovider "github.com/CiscoDevNet/terraform-provider-mso/internal/provider"
	"github.com/CiscoDevNet/terraform-provider-mso/mso"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
)

func main() {
	initLogging(os.Stderr)
	ctx := context.Background()

	server, err := internalprovider.MuxProviderServer(ctx, "dev", mso.Provider().GRPCProvider)
	if err != nil {
		log.Fatalf("[ERROR] Unable to create MSO provider server: %s", err)
	}

	if err := tf6server.Serve(
		"registry.terraform.io/CiscoDevNet/mso",
		server,
	); err != nil {
		log.Fatalf("[ERROR] Unable to serve MSO provider: %s", err)
	}
}

func initLogging(output io.Writer) {
	// tf6server.Serve does not install the standard logger adapter from SDKv2's
	// plugin.Serve. Bridge legacy resource and API client log.Printf calls to
	// structured logs so Terraform can filter them by severity.
	// New Framework resources use tflog with the context passed to their methods;
	// the Framework initializes that contextual logger independently of this adapter.
	logger := hclog.New(&hclog.LoggerOptions{
		Output:     output,
		Level:      hclog.Trace, // Forward all legacy levels; Terraform applies the user's log filter.
		JSONFormat: true,
	})

	log.SetOutput(logger.StandardWriter(&hclog.StandardLoggerOptions{
		InferLevels:              true,
		InferLevelsWithTimestamp: true, // Recognize [ERROR], [TRACE], etc. after Go's timestamp prefix.
	}))
}
