---
layout: "mso"
page_title: "MSO: mso_remote_location"
sidebar_current: "docs-mso-resource-remote-location"
description: |-
  Manages MSO Remote Location
---

# mso_user #

Manages MSO Remote Location.

The default behaviour of the `mso_remote_location` resource stores no sensitive attributes `password`, `ssh_key`, and `passphrase` into the statefile. A change is always detected for the sensitive attributes provided during the plan execution when `store_in_statefile` is not explicitly set to `true`.

## Example Usage ##

```hcl

# remote location with password authentication
resource "mso_remote_location" "password" {
  name        = "remote_location_password"
  description = "remote location with password authentication"
  protocol    = "scp"
  hostname    = "10.0.0.1"
  path        = "/tmp"
  username    = "admin"
  password    = "password"
}

# remote location with ssh key authentication
resource "mso_remote_location" "ssh" {
  name        = "remote_location_ssh"
  description = "remote location with ssh key authentication"
  protocol    = "scp"
  hostname    = "10.0.0.1"
  path        = "/tmp"
  username    = "admin"
  ssh_key     = "ssh_key"
  passphrase  = "passphrase"
}

# remote location with password authentication that stores sensitive attributes to statefile
resource "mso_remote_location" "password" {
  name               = "remote_location_password"
  description        = "remote location with password authentication"
  protocol           = "scp"
  hostname           = "10.0.0.1"
  path               = "/tmp"
  username           = "admin"
  password           = "password"
  store_in_statefile = true
}

```

## Argument Reference ##

* `name` - (Required) The name of the Remote Location.
* `description` - (Optional) The description of the Remote Location.
* `protocol` - (Required) The protocol used to export to the Remote Location. Allowed values are `scp` or `sftp`.
* `hostname` - (Required) The hostname or ip address of the Remote Location.
* `path` - (Required) The full path to a directory on the Remote Location.
* `port` - (Optional) The port used to connect to the Remote Location.  Default to `22`.
* `username`  - (Required) The username used to log in to the Remote Location.
* `password` - (Optional) The password used to log in to the Remote Location.
* `ssh_key` - (Optional) The private ssh key (PEM format) used to log in to the Remote Location.
* `passphrase` - (Optional) The private ssh key passphrase used to log in to the Remote Location.
* `store_in_statefile` - (Optional) Store sensitive attributes `password`, `ssh_key`, and `passphrase` into the statefile. Default to `false`.

## Attribute Reference ##

The only attribute exported with this resource is `id`. Which is set to the id of user associated.

## Importing ##

An existing MSO User can be [imported][docs-import] into this resource via its Id, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_remote_location.example {remote-location-id}
```