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

resource "mso_schema_site_anp_epg_domain" "vmware_domain_with_name_pre_4_2" {
  schema_id                = mso_schema.schema_1.id
  template_name            = one(mso_schema.schema_1.template).name
  site_id                  = data.mso_site.example_site.id
  anp_name                 = mso_schema_template_anp.anp_1.name
  epg_name                 = mso_schema_template_anp_epg.anp_epg_1.name
  domain_type              = "vmmDomain"
  vmm_domain_type          = "VMware"
  domain_name              = "TEST"
  deploy_immediacy         = "immediate"
  resolution_immediacy     = "immediate"
  vlan_encap_mode          = "static"
  allow_micro_segmentation = true
  switching_mode           = "native"
  switch_type              = "default"
  micro_seg_vlan_type      = "vlan"
  micro_seg_vlan           = 46
  port_encap_vlan_type     = "vlan"
  port_encap_vlan          = 45
}

```

### domain_dn usage ###

```hcl

resource "mso_schema_site_anp_epg_domain" "vmware_domain_domain_dn_pre_4_2" {
  schema_id                = mso_schema.schema_1.id
  template_name            = one(mso_schema.schema_1.template).name
  site_id                  = data.mso_site.example_site.id
  anp_name                 = mso_schema_template_anp.anp_1.name
  epg_name                 = mso_schema_template_anp_epg.anp_epg_1.name
  domain_dn                = "uni/vmmp-VMware/dom-TEST"
  deploy_immediacy         = "immediate"
  resolution_immediacy     = "immediate"
  vlan_encap_mode          = "static"
  allow_micro_segmentation = false
  switching_mode           = "native"
  switch_type              = "default"
  micro_seg_vlan_type      = "vlan"
  micro_seg_vlan           = 46
  port_encap_vlan_type     = "vlan"
  port_encap_vlan          = 45
  enhanced_lag_policy_name = "Lacp"
  enhanced_lag_policy_dn   = "uni/vmmp-VMware/dom-TEST/vswitchpolcont/enlacplagp-Lacp"
}


```

### domain_name used in association with domain_type and vmm_domain_type in version >= 4.2 ###

```hcl

resource "mso_schema_site_anp_epg_domain" "vmware_domain_with_name_4_2_up" {
  schema_id                = mso_schema.schema_1.id
  template_name            = one(mso_schema.schema_1.template).name
  site_id                  = data.mso_site.example_site.id
  anp_name                 = mso_schema_template_anp.anp_1.name
  epg_name                 = mso_schema_template_anp_epg.anp_epg_1.name
  domain_type              = "vmmDomain"
  vmm_domain_type          = "VMware"
  domain_name              = "TEST"
  deploy_immediacy         = "immediate"
  resolution_immediacy     = "immediate"
  vlan_encap_mode          = "static"
  allow_micro_segmentation = true
  switching_mode           = "native"
  switch_type              = "default"
  micro_seg_vlan_type      = "vlan"
  micro_seg_vlan           = 46
  port_encap_vlan_type     = "vlan"
  port_encap_vlan          = 45
  delimiter                = "|"
  binding_type             = "static"
  port_allocation          = "fixed"
  num_ports                = 3
  netflow                  = "disabled"
  allow_promiscuous        = "accept"
  mac_changes              = "reject"
  forged_transmits         = "reject"
  custom_epg_name          = "custom_epg_name_1"
}

```

### domain_dn usage in version >= 4.2 ###

```hcl

resource "mso_schema_site_anp_epg_domain" "vmware_domain_with_id_4_2_up" {
  schema_id                = mso_schema.schema_1.id
  template_name            = one(mso_schema.schema_1.template).name
  site_id                  = data.mso_site.example_site.id
  anp_name                 = mso_schema_template_anp.anp_1.name
  epg_name                 = mso_schema_template_anp_epg.anp_epg_1.name
  domain_dn                = "uni/vmmp-VMware/dom-TEST"
  deploy_immediacy         = "immediate"
  resolution_immediacy     = "immediate"
  vlan_encap_mode          = "static"
  allow_micro_segmentation = true
  switching_mode           = "native"
  switch_type              = "default"
  micro_seg_vlan_type      = "vlan"
  micro_seg_vlan           = 46
  port_encap_vlan_type     = "vlan"
  port_encap_vlan          = 45
  delimiter                = "|"
  binding_type             = "static"
  port_allocation          = "fixed"
  num_ports                = 3
  netflow                  = "disabled"
  allow_promiscuous        = "accept"
  mac_changes              = "reject"
  forged_transmits         = "reject"
  custom_epg_name          = "custom_epg_name_1"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg Domain.
* `template_name` - (Required) Template where Anp Epg Domain to be created.
* `site_id` - (Required) SiteID under which you want to deploy Anp Epg Domain.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group to manage.
* `dn` - (Optional) **Deprecated**. The domain profile name. Use `domain_dn` or `domain_name` in association with `domain_type` and `vmm_domain_type` when it is applicable instead.
* `domain_dn` - (Optional) The dn of domain. This is required when `domain_name` and `domain_type` are not specified.
* `domain_name` - (Optional) The domain profile name. This is required when `domain_dn` is not used. This attribute requires `domain_type` and `vmm_domain_type` (when it is applicable) to be set.
* `domain_type` - (Optional) The type of domain to associate. This is required when `domain_dn` is not used. Choices: [ vmmDomain, l3ExtDomain, l2ExtDomain, physicalDomain, fibreChannelDomain ]
* `vmm_domain_type` - (Optional) The vmm domain type. This is required when `domain_type` is vmmDomain and `domain_dn` is not used. Choices: [ VMware, Microsoft, Redhat ]
* `deploy_immediacy` - (Required) The deployment immediacy of the domain. Choices: [ immediate, lazy ]
* `resolution_immediacy` - (Required) Determines when the policies should be resolved and available. Choices: [ immediate, lazy, pre-provision ]
* `vlan_encap_mode` - (Optional) Which VLAN encap mode to use. This attribute can only be used with VMM Domain association. Choices: [ static, dynamic ]
* `allow_micro_segmentation` - (Optional) Specifies microsegmentation is enabled or not. This attribute can only be used with VMM Domain association.
* `switching_mode` - (Optional) Which switching mode to use with this domain association. This attribute can only be used with VMM Domain association.
* `switch_type` - (Optional) Which switch type to use with this domain association. This attribute can only be used with VMM Domain association.
* `micro_seg_vlan_type` - (Optional) Virtual LAN type for microsegmentation. This attribute can only be used with VMM Domain association.
* `micro_seg_vlan` - (Optional) Virtual LAN for microsegmentation. This attribute can only be used with VMM Domain association.
* `port_encap_vlan_type` - (Optional) Virtual LAN type for port encap. This attribute can only be used with VMM Domain association.
* `port_encap_vlan` - (Optional) Virtual LAN for port encap. This attribute can only be used with VMM Domain association.
* `enhanced_lag_policy_name` - (Optional) EPG enhanced lagpolicy name. This attribute can only be used with VMM Domain association.
* `enhanced_lag_policy_dn` - (Optional) Distinguished name of EPG lagpolicy. This attribute can only be used with VMM Domain association.
* `delimiter` - (Optional) The delimiter of the domain. This attribute can only be used with VMM Domain association. Choices: [ |, ~, !, @, ^, +, = ]
* `binding_type` - (Optional) The binding type of the domain. This is required when version of NDO is 4.2+ and can only be used with VMM Domain association. Choices: [ static, dynamic, none, ephemeral ] 
* `port_allocation` - (Optional) The port allocation of the domain. This is required when `binding_type` is static. This attribute can only be used with VMM Domain association. Choices: [ elastic, fixed ]
* `num_ports` - (Optional) The number of ports for the domain. This attribute can only be used with VMM Domain association.
* `netflow ` - (Optional) The netflow preference of the domain. This is required when version of NDO is 4.2+ and can only be used with VMM Domain association. Choices: [ enabled, disabled ]
* `allow_promiscuous` - (Optional) The allow promiscious setting of the domain. This is required when version of NDO is 4.2+ and can only be used with VMM Domain association. Choices: [ accept, reject ]
* `mac_changes` - (Optional) The mac changes setting of the domain. This is required when version of NDO is 4.2+ and can only be used with VMM Domain association. Choices: [ accept, reject ]
* `forged_transmits` - (Optional) The forged transmits setting of the domainn. This is required when version of NDO is 4.2+ and can only be used with VMM Domain association. Choices: [ accept, reject ]
* `custom_epg_name` - (Optional) The custom epg name of the domain. This attribute can only be used with VMM Domain association.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Site Application Network Profiles Endpoint Groups Domain can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_anp_epg_domain.site_anp_epg_domain {schema_id}/sites/{site_id}-{template_name}/anps/{anp_name}/epgs/{epg_name}/domainAssociations/{domain_dn}
```
