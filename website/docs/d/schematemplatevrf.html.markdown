---
layout: "mso"
page_title: "MSO: mso_schema_template_vrf"
sidebar_current: "docs-mso-data-source-schema_template_vrf"
description: |-
  Data source for MSO Schema Template Vrf
---

# mso_schema_site #

Data source for MSO schema template vrf, to fetch the MSO schema template vrf site details.

## Example Usage ##

```hcl
data "mso_schema_template_vrf" "vrf1" {
   schema_id="${mso_schema.schema1.id}"
   template="Template1"
   name= "vrf98"
 }
}
```

## Argument Reference ##


* `schema_id` - (Required) The schema-id where vrf is associated.
* `name` - (Required) name of the vrf to add.



## Attribute Reference ##

* `template` - (Required) template associated with the vrf.