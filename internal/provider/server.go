package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

// MuxProviderServer builds the protocol 6 server used by production and protocol
// tests. Passing the SDKv2 factory keeps this package independent of the legacy
// provider.
func MuxProviderServer(ctx context.Context, version string, sdkFactory func() tfprotov5.ProviderServer) (func() tfprotov6.ProviderServer, error) {
	upgradedSDKServer, err := tf5to6server.UpgradeServer(ctx, sdkFactory)
	if err != nil {
		return nil, err
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx,
		providerserver.NewProtocol6(New(version)()),
		func() tfprotov6.ProviderServer { return upgradedSDKServer },
	)
	if err != nil {
		return nil, err
	}

	return muxServer.ProviderServer, nil
}
