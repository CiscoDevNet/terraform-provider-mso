---
layout: "mso"
page_title: "MSO: mso_schema_template_vrf"
sidebar_current: "docs-mso-resource-schema_template_vrf"
description: |-
  Manages Resource for Schema Template VRF on Cisco Nexus Dashboard Orchestrator (NDO).
---

# mso_schema_template_vrf #

Manages Resource for Schema Template VRF on Cisco Nexus Dashboard Orchestrator (NDO).

## GUI Information ##

* `Location` - Manage -> Tenant Template -> Application -> VRF

## Example Usage ##

```hcl

resource "mso_schema_template_vrf" "example_vrf" {
  schema_id                     = mso_schema.example_schema.id
  template                      = "example_schema_template"
  name                          = "example_vrf"
  display_name                  = "vrf"
  description                   = "Example VRF description"
  layer3_multicast              = true
  vzany                         = false
  ip_data_plane_learning        = "disabled"
  preferred_group               = true
  site_aware_policy_enforcement = true
  rendezvous_points {
    ip_address                      = "1.1.1.2"
    type                            = "static"
    route_map_policy_multicast_uuid = mso_tenant_policies_route_map_policy_multicast.route_map_policy_multicast.uuid
  }
}

```

## Argument Reference ##

* `schema_id` - (Required) The unique ID of the Schema.
* `name` - (Required) The name of the VRF.
* `template` - (Required) The name of the Template associated with the Schema.
* `display_name` - (Required) The name of the VRF as displayed on the NDO/MSO web interface.
* `description` - (Optional) The description of the VRF.
* `layer3_multicast` - (Optional) Whether to enable L3 multicast.
* `vzany` - (Optional) Whether to enable vzany.
* `ip_data_plane_learning` - (Optional) Whether IP data plane learning is enabled or disabled. Allowed values are `disabled`and `enabled`. Default to `enabled`.
* `preferred_group` - (Optional) Whether to enable preferred Endpoint Group.
* `site_aware_policy_enforcement` - (Optional) Whether to enable site aware policy enforcement mode.
* `rendezvous_points` - (Optional) The list of Rendezvous Points. This attribute is supported in NDO v3.0(1) and higher.
  * `rendezvous_points.ip_address` - (Required) The IP Address of the Rendezvous Point.
  * `rendezvous_points.type` - (Required) The type of the Rendezvous Point. Allowed values are `static`, `fabric` and `unknown`.
  * `rendezvous_points.route_map_policy_multicast_uuid` - (Optional) The UUID of the Route Map Policy for Multicast to be associated with the Rendezvous Point.

## Attribute Reference ##

* `uuid` - The NDO UUID of the Route Map Policy for Multicast.

## Importing ##

An existing MSO Resource Schema Template VRF can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_vrf.vrf1 {schema_id}/template/{template}/vrf/{name}
```

