---
layout: "mso"
page_title: "MSO: mso_service_node_type"
sidebar_current: "docs-mso-resource-service_node_type"
description: |-
  Manages MSO Service Node Type
---

# mso_service_node_type #

Manages MSO Service Node Type

## Example Usage ##

```hcl

resource "mso_service_node_type" "node_type" {
  name         = "tftst"
  display_name = "terrform type"
}

```

## Argument Reference ##

* `name` - (Required) Name of the Service Node Type.
* `display_name` - (Optional) Display name of Service Node Type.

## Attribute Reference ##

The only Attribute exposed for this resource is `id`. Which is set to the id of Service Node Type created.

## Importing ##

An existing MSO Service Node Type can be [imported][docs-import] into this resource via its Id, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_service_node_type.node_type {name}
```