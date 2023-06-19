---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg"
sidebar_current: "docs-mso-data-source-schema_template_anp_epg"
description: |-
  Data source for MSO Schema Template Application Network Profiles Endpoint Groups.
---

# mso_schema_template_anp_epg #

Data source for MSO Schema Template Application Network Profiles Endpoint Groups.

## Example Usage ##

```hcl

data "mso_schema_template_anp_epg" "sepg10" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  anp_name      = "ANP"
  name          = "mso_epg1"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg.
* `template_name` - (Required) Template where Anp Epg to be created.
* `anp_name` - (Required) Name of the Application Network Profiles.
* `name` - (Required) Name of the Endpoint Group to manage.

## Attribute Reference ##

* `bd_name` - (Read-Only) Name of the Bridge Domain to associate with.
* `bd_schema_id` - (Read-Only) The schemaID that defines the referenced BD.
* `bd_template_name` - (Read-Only) The template that defines the referenced BD.
* `vrf_name` - (Read-Only) Name of the Vrf.
* `vrf_schema_id` - (Read-Only) The schemaID that defines the referenced VRF.
* `vrf_template_name` - (Read-Only) The template that defines the referenced VRF.
* `display_name` - (Read-Only) The name as displayed on the MSO web interface.
* `description` - (Read-Only) Description of the Anp Epg.
* `useg_epg` - (Read-Only) Boolean flag to enable or disable whether this is a USEG EPG.
* `intra_epg` - (Read-Only) Whether intra EPG isolation is enforced.
* `intersite_multicast_source` - (Read-Only) Whether intersite multicast source is enabled. Default to false.
* `proxy_arp` - (Read-Only) Whether to enable Proxy ARP or not. (For Forwarding control) Default to false.
* `preferred_group` - (Read-Only) Boolean flag to enable or disable whether this EPG is added to preferred group.
* `epg_type` - (Read-Only) Type of the EPG.
* `access_type` - (Read-Only) Access Type of the EPG.
* `deployment_type` - (Read-Only) Deployment Type of the EPG.
* `service_type` - (Read-Only) Service Type of the EPG.
* `custom_service_type` - (Read-Only) Custom Service Type of the EPG.
