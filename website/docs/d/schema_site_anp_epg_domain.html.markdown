---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_domain"
sidebar_current: "docs-mso-data-source-schema_site_anp_epg_domain"
description: |-
  Data source for MSO Schema Site Application Network Profiles Endpoint Groups Domain.
---

# mso_schema_site_anp_epg_domain #

Data source for MSO Schema Site Application Network Profiles Endpoint Groups Domain.

## Example Usage ##

```hcl
data "mso_schema_site_anp_epg_domain" "anpEpgDomain" {
  schema_id = "5e30c4932c00003e5e0a268e"
  template_name = "Template1"
  site_id = "5c7c95d9510000cf01c1ee3d"
  anp_name = "Cloud-First-ANP"
  epg_name = "DB"
  dn = "uni/vmmp-VMware/dom-S2-s2"
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg Domain.
* `template_name` - (Required) Template where Anp Epg Domain to be created.
* `site_id` - (Required) SiteID under which you want to deploy Anp Epg Domain.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group to manage.
* `dn` - (Required) The domain profile name.

## Attribute Reference ##

* `deploy_immediacy` - (Optional) The deployment immediacy of the domain. choices: [ immediate, lazy ]
* `domain_type` - (Optional) The type of domain to associate. choices: [ vmmDomain, l3ExtDomain, l2ExtDomain, physicalDomain, fibreChannel ]
* `resolution_immediacy` - (Optional) Determines when the policies should be resolved and available. choices: [ immediate, lazy, pre-provision ]
* `vlan_encap_mode` - (Optional) Which VLAN enacap mode to use. This attribute can only be used with vmmDomain domain association. choices: [ static, dynamic ]
* `allow_micro_segmentation` - (Optional) Specifies microsegmentation is enabled or not. This attribute can only be used with vmmDomain domain association.
* `switching_mode` - (Optional) Which switching mode to use with this domain association. This attribute can only be used with vmmDomain domain association.
* `switch_type` - (Optional) Which switch type to use with this domain association. This attribute can only be used with vmmDomain domain association.
* `micro_seg_vlan_type` - (Optional) Virtual LAN type for microsegmentation. This attribute can only be used with vmmDomain domain association.
* `micro_seg_vlan` - (Optional) Virtual LAN for microsegmentation. This attribute can only be used with vmmDomain domain association.
* `port_encap_vlan_type` - (Optional) Virtual LAN type for port encap. This attribute can only be used with vmmDomain domain association.
* `port_encap_vlan` - (Optional) Virtual LAN for port encap. This attribute can only be used with vmmDomain domain association.
* `enhanced_lagpolicy_name` - (Optional) EPG enhanced lagpolicy name. This attribute can only be used with vmmDomain domain association.
* `enhanced_lagpolicy_dn` - (Optional) Distinguished name of EPG lagpolicy. This attribute can only be used with vmmDomain domain association.

