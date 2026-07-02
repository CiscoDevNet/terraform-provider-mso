---
layout: "mso"
page_title: "MSO: mso_role"
sidebar_current: "docs-mso-data-source-role"
description: |-
  Data source for MSO Role.
---

# mso_role #

!> **Deprecated** This data source is deprecated: no longer functional on Nexus Dashboard (ND) 3.2+ / NDO 4.2+ and will be removed in the next major release.

Data source for MSO Role.

## Example Usage ##

```hcl

data "mso_role" "example" {
  name  = "UserManager"
}

```

## Argument Reference ##

* `name` - (Required) The name of the Role.

## Attribute Reference #

* `description` - (Read-Only) The description of the Role.
* `display_name` - (Read-Only) The name of the Role as displayed on the MSO UI.
* `read_permissions` - (Read-Only) The read permissions assigned to the Role.
* `write_permissions` - (Read-Only) The write permissions assigned to the Role.
