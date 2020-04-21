package main

import (
	"github.com/ciscoecosystem/terraform-provider-mso/mso"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return mso.Provider()
		},
	})
}
