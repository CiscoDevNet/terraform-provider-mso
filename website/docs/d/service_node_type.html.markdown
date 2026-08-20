---
layout: "mso"
page_title: "MSO: mso_service_node_type"
sidebar_current: "docs-mso-data-source-service_node_type"
description: |-
  Data Source for MSO Service Node Type.
---

# mso_service_node_type #

!> **Deprecated** This data source remains functional on Nexus Dashboard (ND) 4.3+ / NDO 5.3+, but the corresponding `mso_service_node_type` resource is no longer functional. The data source will be removed once ND 4.2 / NDO 5.2 is no longer supported.

Data Source for MSO Service Node Type.

## Example Usage ##

```hcl

data "mso_service_node_type" "example" {
  name = "tftst"
}

```

## Argument Reference ##

* `name` - (Required) The name of the Service Node Type.

## Attribute Reference ##

* `display_name` - (Read-Only) The name of the Service Node Type as displayed on the MSO UI.
