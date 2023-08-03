---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg_contract"
sidebar_current: "docs-mso-data-source-schema_template_anp_epg_contract"
description: |-
  Data source for MSO Schema Template Application Network Profiles Endpoint Group Contract.
---

# mso_schema_template_anp_epg_contract #

Data source for MSO Schema Template Application Network Profiles Endpoint Group Contract.

## Example Usage ##

```hcl

data "mso_schema_template_anp_epg_contract" "example" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  anp_name      = "WoS-Cloud-Only-2"
  epg_name      = "DB"
  contract_name = "Web2-to-DB2"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the EPG.
* `template_name` - (Required) The template name of the EPG.
* `anp_name` - (Required) The name of the ANP.
* `epg_name` - (Required) The name of the EPG.
* `contract_name` - (Required) The name of the Contract.
* `contract_schema_id` - (Optional) The schema ID of the Contract. The `schema_id` of the EPG. will be used if not provided. 
* `contract_template_name` - (Optional) The template name of the Contract. The `contract_template_name` of the EPG. will be used if not provided. 

## Attribute Reference ##

* `relationship_type` - (Read-Only) The relationship type of the Contract.
