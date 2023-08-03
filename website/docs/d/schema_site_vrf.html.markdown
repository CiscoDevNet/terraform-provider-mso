---
layout: "mso"
page_title: "MSO: mso_schema_site_vrf"
sidebar_current: "docs-mso-data-source-schema_site_vrf"
description: |-
  Data source for MSO Schema Site VRF.
---

# mso_schema_site_vrf #

Data source for MSO Schema Site VRF.

## Example Usage ##

```hcl

data "mso_schema_site_vrf" "example" {
  schema_id     = data.mso_schema.schema1.id
  site_id       = data.mso_site.site1.id
  template_name = "Template1"
  vrf_name      = "vrf5810"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the VRF is deployed.
* `site_id` - (Required) The site ID under which the VRF is deployed.
* `template_name` - (Required) The template name under which the VRF is deployed.
* `vrf_name` - (Required) The name of the VRF.
