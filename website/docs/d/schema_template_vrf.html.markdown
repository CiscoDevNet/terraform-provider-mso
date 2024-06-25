---
layout: "mso"
page_title: "MSO: mso_schema_template_vrf"
sidebar_current: "docs-mso-data-source-schema_template_vrf"
description: |-
  Data source for MSO Schema Template VRF.
---

# mso_schema_template_vrf #

Data source for MSO Schema Template VRF.

## Example Usage ##

```hcl

data "mso_schema_template_vrf" "example" {
  schema_id = data.mso_schema.schema1.id
  template  = "Template1"
  name      = "vrf98"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the VRF.
* `template` - (Required) The template name of the VRF.
* `name` - (Required) The name of the VRF.

## Attribute Reference ##

* `display_name` - (Read-Only) The name of the VRF as displayed on the MSO UI.
* `layer3_multicast` - (Read-Only) Whether L3 multicast is enabled.
* `vzany` - (Read-Only) Whether vzany is enabled.
* `ip_data_plane_learning` - (Read-Only) Whether IP data plane learning is enabled.
* `preferred_group` - (Read-Only) Whether to preferred group is enabled.
* `description` - (Read-Only) The description of the VRF.
* `site_aware_policy_enforcement` - (Read-Only) Whether site aware policy enforcement mode is enabled.
