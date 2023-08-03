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

* `schema_id` - (Required) The schema ID of the Contract.
* `template_name` - (Required) The template name of the Contract.
* `anp_name` - (Required) The name of the ANP.
* `epg_name` - (Required) The name of the EPG.
* `contract_name` - (Required) The name of the Contract.

## Attribute Reference ##
* `contract_schema_id` - (Read-Only) The schema ID associated with the Contract.
* `contract_template_name` - (Read-Only) The template name associated with the Contract.
* `relationship_type` - (Read-Only) The relationship type of the Contract.
