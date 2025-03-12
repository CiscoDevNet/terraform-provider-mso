---
layout: "mso"
page_title: "MSO: mso_schema_template_vrf"
sidebar_current: "docs-mso-data-source-schema_template_vrf"
description: |-
  Data source for MSO Schema Template VRF.
---

# mso_schema_template_vrf #

Data source for MSO Schema Template VRF.

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Application -> VRF

## Example Usage ##

```hcl

data "mso_schema_template_vrf" "example_vrf" {
  schema_id = data.mso_schema.example_schema.id
  template  = "example_schema_template"
  name      = "example_vrf"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the VRF.
* `template` - (Required) The name of the Template associated to the Schema.
* `name` - (Required) The name of the VRF.

## Attribute Reference ##

* `uuid` - (Read-Only) The UUID of the VRF.
* `display_name` - (Read-Only) The name of the VRF as displayed on the MSO UI.
* `layer3_multicast` - (Read-Only) Whether L3 multicast is enabled.
* `vzany` - (Read-Only) Whether vzany is enabled.
* `ip_data_plane_learning` - (Read-Only) Whether IP data plane learning is enabled.
* `preferred_group` - (Read-Only) Whether to preferred group is enabled.
* `description` - (Read-Only) The description of the VRF.
* `site_aware_policy_enforcement` - (Read-Only) Whether site aware policy enforcement mode is enabled.
* `rendezvous_points` - (Read-Only) The list of Rendezvous Points.
  * `rendezvous_points.ip_address` - (Read-Only) The IP Address of the Rendezvous Point.
  * `rendezvous_points.type` - (Read-Only) The type of the Rendezvous Point.  Allowed values are `static`, `fabric` and `unknown`.
  * `rendezvous_points.mutlicast_route_map_policy_uuid` - (Read-Only) The UUID of the Route Map Policy for Multicast associated with the Rendezvous Point.