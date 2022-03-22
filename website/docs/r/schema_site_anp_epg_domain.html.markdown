---
layout: "mso"
page_title: "MSO: mso_schema_site_anp_epg_domain"
sidebar_current: "docs-mso-resource-schema_site_anp_epg_domain"
description: |-
  Manages MSO Schema Site Application Network Profiles Endpoint Groups Domain.
---

# mso_schema_site_anp_epg_domain #

Manages MSO Schema Site Application Network Profiles Endpoint Groups Domain.

## Example Usage ##

### domain_name used in association with domain_type and vmm_domain_type ###

```hcl
resource "mso_schema_site_anp_epg_domain" "site_anp_epg_domain" {
  schema_id = "5c4d9fca270000a101f8094a"
  template_name = "Template1"
  site_id = "5c7c95b25100008f01c1ee3c"
  anp_name = "ANP"
  epg_name = "Web"
  domain_type = "vmmDomain"
  vmm_domain_type = "VMware"
  domain_name = "VMware-ab"
  deploy_immediacy = "lazy"
  resolution_immediacy = "lazy"
  vlan_encap_mode = "static"
  allow_micro_segmentation = true
  switching_mode = "native"
  switch_type = "default"
  micro_seg_vlan_type = "vlan"
  micro_seg_vlan = 46
  port_encap_vlan_type = "vlan"
  port_encap_vlan = 45
  enhanced_lag_policy_name = "name"
  enhanced_lag_policy_dn = "dn"

}

```

### domain_dn usage ###

```hcl
resource "mso_schema_site_anp_epg_domain" "site_anp_epg_domain" {
  schema_id = "5c4d9fca270000a101f8094a"
  template_name = "Template1"
  site_id = "5c7c95b25100008f01c1ee3c"
  anp_name = "ANP"
  epg_name = "Web"
  domain_dn = "uni/vmmp-VMware/dom-VMware-ab"
  deploy_immediacy = "lazy"
  resolution_immediacy = "lazy"
  vlan_encap_mode = "static"
  allow_micro_segmentation = true
  switching_mode = "native"
  switch_type = "default"
  micro_seg_vlan_type = "vlan"
  micro_seg_vlan = 46
  port_encap_vlan_type = "vlan"
  port_encap_vlan = 45
  enhanced_lag_policy_name = "name"
  enhanced_lag_policy_dn = "dn"

}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg Domain.
* `template_name` - (Required) Template where Anp Epg Domain to be created.
* `site_id` - (Required) SiteID under which you want to deploy Anp Epg Domain.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group to manage.
* `dn` - (Optional) **Deprecated**. The domain profile name. This is required when `domain_dn` is not used. This attribute requires `domain_type` and `vmm_domain_type` (if applicable) to be set.
* `domain_name` - (Optional) The domain profile name. This is required when `domain_dn` is not used. This attribute requires `domain_type` and `vmm_domain_type` (if applicable) to be set.
* `domain_dn` - (Optional) The dn of domain. This is required when `domain_name` and `domain_type` are not specified.
* `vmm_domain_type` - (Required) The vmm domain type. This is required when `domain_type` is vmmDomain. choices: [ VMware, Microsoft, Redhat ]
* `deploy_immediacy` - (Required) The deployment immediacy of the domain. choices: [ immediate, lazy ]
* `domain_type` - (Required) The type of domain to associate. choices: [ vmmDomain, l3ExtDomain, l2ExtDomain, physicalDomain, fibreChannelDomain ]
* `resolution_immediacy` - (Required) Determines when the policies should be resolved and available. choices: [ immediate, lazy, pre-provision ]
* `vlan_encap_mode` - (Optional) Which VLAN encap mode to use. This attribute can only be used with vmmDomain domain association. choices: [ static, dynamic ]
* `allow_micro_segmentation` - (Optional) Specifies microsegmentation is enabled or not. This attribute can only be used with vmmDomain domain association.
* `switching_mode` - (Optional) Which switching mode to use with this domain association. This attribute can only be used with vmmDomain domain association.
* `switch_type` - (Optional) Which switch type to use with this domain association. This attribute can only be used with vmmDomain domain association.
* `micro_seg_vlan_type` - (Optional) Virtual LAN type for microsegmentation. This attribute can only be used with vmmDomain domain association.
* `micro_seg_vlan` - (Optional) Virtual LAN for microsegmentation. This attribute can only be used with vmmDomain domain association.
* `port_encap_vlan_type` - (Optional) Virtual LAN type for port encap. This attribute can only be used with vmmDomain domain association.
* `port_encap_vlan` - (Optional) Virtual LAN for port encap. This attribute can only be used with vmmDomain domain association.
* `enhanced_lag_policy_name` - (Optional) EPG enhanced lagpolicy name. This attribute can only be used with vmmDomain domain association.
* `enhanced_lag_policy_dn` - (Optional) Distinguished name of EPG lagpolicy. This attribute can only be used with vmmDomain domain association.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Site Application Network Profiles Endpoint Groups Domain can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_anp_epg_domain.site_anp_epg_domain {schema_id}/sites/{site_id}-{template_name}/anps/{anp_name}/epgs/{epg_name}/domainDn/{domain_dn}
```
