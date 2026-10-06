package main

import (
	"context"
	"log"

	internalprovider "github.com/CiscoDevNet/terraform-provider-mso/internal/provider"
	"github.com/CiscoDevNet/terraform-provider-mso/mso"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func main() {
	ctx := context.Background()

	server, err := muxProviderServer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := tf6server.Serve(
		"registry.terraform.io/CiscoDevNet/mso",
		server,
	); err != nil {
		log.Fatal(err)
	}
}

func muxProviderServer(ctx context.Context) (func() tfprotov6.ProviderServer, error) {
	upgradedSDKServer, err := tf5to6server.UpgradeServer(ctx, mso.Provider().GRPCProvider)
	if err != nil {
		return nil, err
	}

	providers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(internalprovider.New("dev")()),
		func() tfprotov6.ProviderServer {
			return upgradedSDKServer
		},
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		return nil, err
	}

	return muxServer.ProviderServer, nil
}
