---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg_contract"
sidebar_current: "docs-mso-data-source-schema_template_anp_epg_contract"
description: |-
  MSO Schema Template Application Network Profile(ANP) Endpoint Group(EPG) Contract Data Source
---

# mso_schema_template_anp_epg_contract #

MSO Schema Template ANP EPG Contract Data source.

## Example Usage ##

```hcl

data "mso_schema_template_anp_epg_contract" "contract" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  anp_name      = "WoS-Cloud-Only-2"
  epg_name      = "DB"
  contract_name = "Web2-to-DB2"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID of Application Network Profile(ANP) EPG.
* `template_name` - (Required) Template name of ANP EPG .
* `anp_name` - (Required) Name of the Application Network Profile.
* `epg_name` - (Required) Name of the Endpoint Group.
* `relationship_type` - (Optional) Relationship Type of the ANP EPG Contract on MSO UI.
* `contract_name` - (Required) Name of the contract.



## Attribute Reference ##
* `contract_schema_id` - (Optional) SchemaID associated with the contract.
* `contract_template_name` - (Optional) Template Name associated with the contract.

