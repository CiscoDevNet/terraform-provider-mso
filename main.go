package main

import (
	"github.com/CiscoDevNet/terraform-provider-mso/mso"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return mso.Provider()
		},
	})
}
