---
layout: "mso"
page_title: "MSO: mso_schema_template_external_epg_contract"
sidebar_current: "docs-mso-data-source-schema_template_external_epg_contract"
description: |-
  MSO Schema Template External End Point Group Contract Data source.
---

# mso_schema_template_external_epg_contract #

MSO Schema Template External End Point Group Contract Data source.

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

## Attribute Reference ##

* `relationship_type` - (Read-Only) The relationship type of the Contract.
* `contract_schema_id` - (Read-Only) The schema ID of the Contract.
* `contract_template_name` - (Read-Only) The template name of the Contract.
