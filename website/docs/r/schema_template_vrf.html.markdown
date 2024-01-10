---
layout: "mso"
page_title: "MSO: mso_schema_template_vrf"
sidebar_current: "docs-mso-resource-schema_template_vrf"
description: |-
  Manages MSO Resource Schema Template Vrf
---

# mso_schema_template_vrf #

Manages MSO Resource Schema Template Vrf

## Example Usage ##

```hcl

resource "mso_schema_template_vrf" "vrf1" {
  schema_id        = mso_schema.schema1.id
  template         = "temp3"
  name             = "vrf982"
  display_name     = "vz1"
  layer3_multicast = false
  vzany            = false
  ip_data_plane_learning = "disabled"
  preferred_group = true
}

```

## Argument Reference ##


* `schema_id` - (Required) The schema-id where vrf is associated.
* `name` - (Required) name of the vrf to add.
* `template` - (Required) template associated with the vrf.
* `display_name` - (Required) The name as displayed on the MSO web interface.
* `description` - (Optional) The description of the vrf.
* `layer3_multicast` - (Optional) Whether to enable L3 multicast.
* `vzany` - (Optional) Whether to enable vzany.
* `ip_data_plane_learning` - (Optional) Whether IP data plane learning is enabled or disabled. Allowed values are `disabled`and `enabled`. Default to `enabled`.
* `preferred_group` - (Optional) Whether to enable preferred Endpoint Group.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Resource Schema Template Vrf can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_vrf.vrf1 {schema_id}/template/{template}/vrf/{name}
```

