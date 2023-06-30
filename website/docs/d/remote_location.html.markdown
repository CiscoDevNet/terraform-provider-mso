---
layout: "mso"
page_title: "MSO: mso_remote_location"
sidebar_current: "docs-mso-resource-remote-location"
description: |-
  Data source for MSO Remote Location
---

# mso_user #

Data source for MSO Remote Location

## Example Usage ##

```hcl

data "mso_remote_location" "password" {
  name = "remote_location_name"
}

```

## Argument Reference ##

* `name` - (Required) The name of the Remote Location.

## Attribute Reference ##

* `description` - (Read-Only) The description of the Remote Location.
* `protocol` - (Read-Only) The protocol used to export to the Remote Location.
* `hostname` - (Read-Only) The hostname or ip address of the Remote Location.
* `path` - (Read-Only) The full path to a directory on the Remote Location.
* `port` - (Read-Only) The port used to connect to the Remote Location.
* `username`  - (Read-Only) The username used to log in to the Remote Location.
