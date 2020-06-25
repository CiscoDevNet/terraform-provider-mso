---
layout: "mso"
page_title: "MSO: mso_schema_template_vrf"
sidebar_current: "docs-mso-data-source-schema_template_vrf"
description: |-
  Data source for MSO Schema Template Vrf
---

# mso_schema_template_vrf #

Data source for MSO schema template vrf, to fetch the MSO schema template vrf site details.

## Example Usage ##

```hcl
data "mso_schema_template_vrf" "vrf1" {
  schema_id = "${mso_schema.schema1.id}"
  template  = "Template1"
  name      = "vrf98"
}
```

## Argument Reference ##

* `schema_id` - (Required) The schema-id where vrf is associated.
* `name` - (Required) name of the vrf to add.
* `template` - (Required) template associated with the vrf.

## Attribute Reference ##

* `display_name` - (Optional) The name as displayed on the MSO web interface.
* `layer3_multicast` - (Optional) Whether to enable L3 multicast.
* `vzany` - (Optional) Whether to enable vzany.

