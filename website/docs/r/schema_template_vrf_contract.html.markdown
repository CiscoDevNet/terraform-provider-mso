---
layout: "mso"
page_title: "MSO: mso_schema_template_vrf_contract"
sidebar_current: "docs-mso-resource-schema_template_vrf_contract"
description: |-
  Manages MSO Resource Schema Template Vrf Contract.
---

# mso_schema_template_vrf_contract #

Manages MSO Resource Schema Template Vrf Contract. This resource is supported in MSO v3.0 or higher.

## Example Usage ##

```hcl

resource "mso_schema_template_vrf_contract" "demovrf01" {
  schema_id              = mso_schema.schema1.id
  template_name          = "Template1"
  vrf_name               = "myVrf"
  relationship_type      = "provider"
  contract_name          = "hubcon"
  contract_schema_id     = "5efd6ea60f00005b0ebbd643"
  contract_template_name = "Template1"
}

```

## Argument Reference ##


* `schema_id` - (Required) The schema-id where vrf is associated.
* `template_name` - (Required) template associated with the vrf.
* `vrf_name` - (Required) name of the vrf with contract to be attached.
* `relationship_type` - (Required) Type of relation between VRF and Contract. Allowed values are `provider` and `consumer`. Provider contracts cannot be shared across the VRF.
* `contract_name` - (Required) Name of contract to be attached with the VRF.
* `contract_schema_id` - (Optional) SchemaId of contract. This parameter should be used when the contract and VRF are in different schemas. `schema_id` will be used if not provided.
* `contract_template_name` - (Optional) Name of template where contract is residing. This parameter should be used when the contract and VRF are in different Templates. `template_name` will be used if not provided.


## Attribute Reference ##
The only attribute exported is `id`. Which is set to the name of contract attached.

## Importing ##

An existing MSO Resource Schema Template Vrf Contract can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_vrf_contract.demovrf01 {schema_id}/template/{template_name}/vrf/{vrf_name}/contract/{contract_name}/type/{relationship_type}
```