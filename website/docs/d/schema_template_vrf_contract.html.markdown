---
layout: "mso"
page_title: "MSO: mso_schema_template_vrf_contract"
sidebar_current: "docs-mso-data-source-schema_template_vrf_contract"
description: |-
  Data Source for MSO Schema Template Vrf Contract.
---

# mso_schema_template_vrf_contract #

Data Source for MSO  Schema Template Vrf Contract. This data source is supported in MSO v3.0 or higher.

## Example Usage ##

```hcl

data "mso_schema_template_vrf_contract" "demovrf01" {
  schema_id         = data.mso_schema.schema1.id
  template_name     = "Template1"
  vrf_name          = "myVrf"
  relationship_type = "provider"
  contract_name     = "hubcon"
}

```

## Argument Reference ##


* `schema_id` - (Required) The schema-id where vrf is associated.
* `template_name` - (Required) template associated with the vrf.
* `vrf_name` - (Required) name of the vrf with contract to be attached.
* `relationship_type` - (Required) Type of relation between VRF and Contract. Allowed values are `provider` and `consumer`.
* `contract_name` - (Required) Name of contract to be attached with the VRF.



## Attribute Reference ##
* `contract_schema_id` - (Optional) SchemaId of contract. 
* `contract_template_name` - (Optional) Name of template where contract is residing.


