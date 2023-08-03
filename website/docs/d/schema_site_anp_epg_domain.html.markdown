---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_domain"
sidebar_current: "docs-mso-data-source-schema_site_anp_epg_domain"
description: |-
  Data source for MSO Schema Site Application Network Profiles End Point Group Domain.
---

# mso_schema_site_anp_epg_domain #

Data source for MSO Schema Site Application Network Profiles End Point Group Domain.

## Example Usage ##

### domain_name used in association with domain_type and vmm_domain_type ###

```hcl

data "mso_schema_site_anp_epg_domain" "example_name" {
  schema_id       = data.mso_schema.schema1.id
  site_id         = data.mso_site.site1.id
  template_name   = "Template1"
  anp_name        = "ANP"
  epg_name        = "Web"
  domain_name     = "VMware-ab"
  domain_type     = "vmmDomain"
  vmm_domain_type = "VMware"
}

```

### domain_dn usage ###

```hcl

resource "mso_schema_site_anp_epg_domain" "example_dn" {
  schema_id = "5c4d9fca270000a101f8094a"
  template_name = "Template1"
  site_id = "5c7c95b25100008f01c1ee3c"
  anp_name = "ANP"
  epg_name = "Web"
  domain_dn = "uni/vmmp-VMware/dom-VMware-ab"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID under which the Domain is deployed.
* `site_id` - (Required) The site ID under which the Domain is deployed.
* `template_name` - (Required) The template name under which the Domain is deployed.
* `anp_name` - (Required) The ANP name under which the Domain is deployed.
* `epg_name` - (Required) The EPG name under which the Domain is deployed.
* `domain_dn` - (Optional) The DN of the Domain. This is required when `domain_name` and `domain_type` are not specified.
* `domain_name` - (Optional) The name of the Domain. This is required when `domain_dn` is not used. This attribute requires `domain_type` and `vmm_domain_type` (when it is applicable) to be set.
* `domain_type` - (Optional)  The type of the Domain. This is required when `domain_dn` is not used. Choices: [ vmmDomain, l3ExtDomain, l2ExtDomain, physicalDomain, fibreChannelDomain ]
* `vmm_domain_type` - (Optional) The type of the VMM Domain. This is required when `domain_type` is vmmDomain and `domain_dn` is not used. Choices: [ VMware, Microsoft, Redhat ]

## Attribute Reference ##

* `template_name` - (Read-Only) The template of the Domain.
* `deploy_immediacy` - (Read-Only) The deployment immediacy of the Domain.
* `resolution_immediacy` - (Read-Only) The resolution immediacy of the Domain.
* `vlan_encap_mode` - (Read-Only) The VLAN encapsulation mode of the Domain.
* `allow_micro_segmentation` - (Read-Only) The allow microsegmentation of the Domain.
* `switching_mode` - (Read-Only) The switching mode of the Domain. 
* `switch_type` - (Read-Only) The switch type of the Domain.
* `micro_seg_vlan_type` - (Read-Only) The virtual LAN type for microsegmentation of the Domain. 
* `micro_seg_vlan` - (Read-Only) The virtual LAN for microsegmentation of the Domain. 
* `port_encap_vlan_type` - (Read-Only) The virtual LAN type for port encapsulation of the Domain.
* `port_encap_vlan` - (Read-Only) The port encapapsulation of the Domain.
* `enhanced_lag_policy_name` - (Read-Only) The EPG enhanced lag policy name of the Domain.
* `enhanced_lag_policy_dn` - (Read-Only) The EPG enhanced lag policy DN of the Domain. 
