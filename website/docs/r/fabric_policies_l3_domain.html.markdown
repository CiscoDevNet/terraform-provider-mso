---
layout: "mso"
page_title: "MSO: mso_fabric_policies_l3_domain"
sidebar_current: "docs-mso-resource-fabric_policies_l3_domain"
description: |-
  Manages L3 Domains on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_fabric_policies_l3_domain #

Manages Layer 3 (L3) Domains on Cisco Nexus Dashboard Orchestrator (NDO). This resource is only supported NDO v4.3 and later.

## GUI Information ##

* `Location` - Manage -> Fabric Template -> Fabric Policies -> L3 Domains

## Example Usage ##

```hcl
resource "mso_fabric_policies_l3_domain" "l3out_domain" {
  template_id    = mso_template.fabric_template.id
  name           = "l3out_domain"
  description    = "L3 Domain for external routing"
  vlan_pool_uuid = mso_fabric_policies_vlan_pool.l3out_pool.uuid
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the fabric policy template.
* `name` - (Required) The name of the L3 Domain.
* `description` - (Optional) The description of the L3 Domain.
* `vlan_pool_uuid` - (Optional) The NDO UUID of the VLAN Pool to associate with this L3 Domain. 

## Attribute Reference ##

* `uuid` - The NDO UUID of the L3 Domain.
* `id` - The unique Terraform identifier of the L3 Domain in the template.

## Importing ##

An existing MSO L3 Domain can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: https://www.terraform.io/docs/import/index.html


```bash
terraform import mso_fabric_policies_l3_domain.l3_domain templateId/{template_id}/L3Domain/{name}
```