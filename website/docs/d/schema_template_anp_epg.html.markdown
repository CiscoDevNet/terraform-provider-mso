---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg"
sidebar_current: "docs-mso-data-source-schema_template_anp_epg"
description: |-
  Data source for MSO Schema Template Application Network Profiles Endpoint Group (EPG).
---

# mso_schema_template_anp_epg #

Data source for MSO Schema Template Application Network Profiles Endpoint Group (EPG).

## Example Usage ##

```hcl

data "mso_schema_template_anp_epg" "example" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  anp_name      = "ANP"
  name          = "mso_epg1"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the EPG.
* `template_name` - (Required) The template name of the EPG.
* `anp_name` - (Required) The name of the ANP.
* `name` - (Required) The name of the EPG.

## Attribute Reference ##

* `uuid` - (Read-Only) The UUID of the EPG.
* `bd_name` - (Read-Only) The name of the BD associated with the EPG.
* `bd_schema_id` - (Read-Only) The schema ID of the BD associated with the EPG.
* `bd_template_name` - (Read-Only) The template name of the BD associated with the EPG.
* `vrf_name` - (Read-Only) The name of the VRF associated with the EPG.
* `vrf_schema_id` - (Read-Only) The schema ID of the VRF associated with the EPG.
* `vrf_template_name` - (Read-Only) The template name of the VRF associated with the EPG.
* `display_name` - (Read-Only) The name of the EPG as displayed on the MSO UI.
* `description` - (Read-Only) The description of the EPG.
* `useg_epg` - (Read-Only) Whether the EPG is a uSeg EPG.
* `intra_epg` - (Read-Only) Whether intra EPG isolation is enforced.
* `intersite_multicast_source` - (Read-Only) Whether intersite multicast source is enabled.
* `proxy_arp` - (Read-Only) Whether Proxy ARP is enabled.
* `preferred_group` - (Read-Only)  Whether the EPG is added to preferred group.
* `epg_type` - (Read-Only) The type of the EPG.
* `access_type` - (Read-Only) The access type of the EPG.
* `deployment_type` - (Read-Only) The deployment type of the EPG.
* `service_type` - (Read-Only) The service type of the EPG.
* `custom_service_type` - (Read-Only) The custom service type of the EPG.
