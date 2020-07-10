---
layout: "mso"
page_title: "MSO: mso_service_node_type"
sidebar_current: "docs-mso-data-source-service_node_type"
description: |-
  Data Source for MSO Service Node Type
---

# mso_service_node_type #

Data Source for MSO Service Node Type

## Example Usage ##

```hcl
data "mso_service_node_type" "node_data" {
  name = "tftst"
}

```

## Argument Reference ##

* `name` - (Required) Name of the Service Node Type.


## Attribute Reference ##

* `display_name` - (Optional) Display name of Service Node Type.