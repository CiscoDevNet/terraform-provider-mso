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
* `anp_name` - (Required) Name of Application Network Profiles.
* `name` - (Required) Name of Endpoint Group to manage.

## Attribute Reference ##

* `bd_name` - (Optional) Name of Bridge Domain to associate with.
* `bd_schema_id` - (Opional) The schemaID that defines the referenced BD.
* `bd_template_name` - (Optional) The template that defines the referenced BD.
* `vrf_name` - (Optional) Name of Vrf.
* `vrf_schema_id` - (Optional) The schemaID that defines the referenced VRF.
* `vrf_template_name` - (Optional) The template that defines the referenced VRF.
* `display_name` - (Optional) The name as displayed on the MSO web interface.
* `useg_epg` - (Optional) Boolean flag to enable or disable whether this is a USEG EPG.
* `intra_epg` - (Optional) Whether intra EPG isolation is enforced. choices: [ enforced, unenforced ]
* `intersite_multicast_source` - (Optional) Whether intersite multicast source is enabled. Default to false.
* `proxy_arp` - (Optional) Whether to enable Proxy ARP or not. (For Forwarding control) Default to false.
* `preferred_group` - (Optional) Boolean flag to enable or disable whether this EPG is added to preferred group.