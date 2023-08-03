---
layout: "mso"
page_title: "MSO: mso_schema_template_vrf_contract"
sidebar_current: "docs-mso-data-source-schema_template_vrf_contract"
description: |-
  Data Source for MSO Schema Template VRF Contract.
---

# mso_schema_template_vrf_contract #

Data Source for MSO Schema Template VRF Contract. This data source is supported in MSO v3.0 or higher.

## Example Usage ##

```hcl

data "mso_schema_template_vrf_contract" "example" {
  schema_id         = data.mso_schema.schema1.id
  template_name     = "Template1"
  vrf_name          = "myVrf"
  relationship_type = "provider"
  contract_name     = "hubcon"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the VRF.
* `template_name` - (Required) The template name of the VRF.
* `vrf_name` - (Required) The name of the VRF.
* `relationship_type` - (Required) The relationship type of the VRF with Contract. Allowed values are `provider` and `consumer`.
* `contract_name` - (Required) The name of the Contract.
* `contract_schema_id` - (Optional) The schema ID of the Contract. The `schema_id` of the VRF will be used if not provided. 
* `contract_template_name` - (Optional) The template name of the Contract. The `contract_template_name` of the VRF will be used if not provided. 
