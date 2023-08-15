---
layout: "mso"
page_title: "MSO: mso_schema_template_external_epg_contract"
sidebar_current: "docs-mso-data-source-schema_template_external_epg_contract"
description: |-
  Data source for MSO Schema Template External End Point Group Contract.
---

# mso_schema_template_external_epg_contract #

Data source for MSO Schema Template External End Point Group Contract.

## Example Usage ##

```hcl

data "mso_schema_template_external_epg_contract" "example" {
  schema_id         = data.mso_schema.schema1.id
  template_name     = "Template1"
  external_epg_name = "ExternalEPG1"
  contract_name     = "contract1006"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the External EPG.
* `template_name` - (Required) The template name of the External EPG.
* `external_epg_name` - (Required) The name of the External EPG.
* `contract_name` - (Required) The name of the Contract.
* `contract_schema_id` - (Optional) The schema ID of the Contract. The `schema_id` of the External EPG. will be used if not provided. 
* `contract_template_name` - (Optional) The template name of the Contract. The `contract_template_name` of the External EPG. will be used if not provided. 

## Attribute Reference ##

* `relationship_type` - (Read-Only) The relationship type of the Contract.
