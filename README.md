# Cisco MSO Provider

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) v0.12 or later.

- [Go](https://golang.org/doc/install) Latest Version

## Building The Provider ##
Clone this repository to: `$GOPATH/src/github.com/CiscoDevNet/terraform-provider-mso`.
Clone mso-go-client to `$GOPATH/src/github.com/ciscoecosystem/mso-go-client`.

```sh
$ mkdir -p $GOPATH/src/github.com/ciscoecosystem; cd $GOPATH/src/github.com/ciscoecosystem
$ git clone https://github.com/CiscoDevNet/terraform-provider-mso.git
$ git clone https://github.com/ciscoecosystem/mso-go-client.git
```


Using The Provider
------------------
If you are building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory, run `terraform init` to initialize it.

ex.
```hcl
terraform {
  required_providers {
    mso = {
      source = "CiscoDevNet/mso"
    }
  }
}
#configure provider with your cisco mso credentials.
provider "mso" {
  # cisco-mso user name
  username = "admin"
  # cisco-mso password
  password = "password"
  # cisco-mso url
  url      = "https://my-cisco-mso.com"
  insecure = true
  proxy_url = "https://proxy_server:proxy_port"
  platform = "nd"
}

resource "mso_schema" "schema1" {
  name          = "nkp1002"
  template_name = "temp1"
  tenant_id     = "5e9d09482c000068500a269a"

}

```


Developing The Provider
-----------------------
If you want to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine. You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider with sanity checks present in scripts directory and put the provider binary in `$GOPATH/bin` directory.

