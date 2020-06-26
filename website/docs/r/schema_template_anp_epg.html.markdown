---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg"
sidebar_current: "docs-mso-resource-schema_template_anp_epg"
description: |-
  Manages MSO Schema Template Application Network Profiles Endpoint Groups.
---

# mso_schema_template_anp_epg #

Manages MSO Schema Template Application Network Profiles Endpoint Groups.

## Example Usage ##

```hcl
resource "mso_schema_template_anp_epg" "anp_epg" {
  schema_id = "5c4d5bb72700000401f80948"
  template_name = "Template1"
  anp_name = "ANP"
  name = "mso_epg1"
  bd_name = "BD1"
  vrf_name = "DEVNET-VRF"
  display_name = "mso_epg1"
  useg_epg = true
  intra_epg = "unenforced"
  intersite_multicast_source = false
  preferred_group = false
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg.
* `template_name` - (Required) Template where Anp Epg to be created.
* `anp_name` - (Required) Name of Application Network Profiles.
* `name` - (Required) Name of Endpoint Group to manage.
* `bd_name` - (Required) Name of Bridge Domain to associate with.
* `bd_schema_id` - (Opional) The schemaID that defines the referenced BD.
* `bd_template_name` - (Optional) The template that defines the referenced BD.
* `vrf_name` - (Required) Name of Vrf.
* `vrf_schema_id` - (Optional) The schemaID that defines the referenced VRF.
* `vrf_template_name` - (Optional) The template that defines the referenced VRF.
* `display_name` - (Optional) The name as displayed on the MSO web interface.
* `useg_epg` - (Optional) Boolean flag to enable or disable whether this is a USEG EPG. Default value is set to false.
* `intra_epg` - (Optional) Whether intra EPG isolation is enforced. choices: [ enforced, unenforced ]
* `intersite_multicast_source` - (Optional) Whether intersite multicast source is enabled. Default to false.
* `proxy_arp` - (Optional) Whether to enable Proxy ARP or not. (For Forwarding control) Default to false.
* `preferred_group` - (Optional) Boolean flag to enable or disable whether this EPG is added to preferred group.      Default value is set to false.

## Attribute Reference ##

No attributes are exported.
