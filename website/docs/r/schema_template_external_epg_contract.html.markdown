---
layout: "mso"
page_title: "MSO: mso_schema_template_external_epg_contract"
sidebar_current: "docs-mso-resource-schema_template_external_epg_contract"
description: |-
  Manages MSO Schema Template External Endpoint Group Contract.
---

# mso_schema_template_external_epg_contract #

Manages MSO Schema Template External Endpoint Group Contract.

## Example Usage ##

```hcl
resource "mso_schema_template_external_epg_contract" "c1" {
  schema_id = "5ea809672c00003bc40a2799"
  template_name = "Template1"
  contract_name = "contract1006"
  external_epg_name = "UntitledExternalEPG1"
  relationship_type = "consumer"
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy External-epg.
* `template_name` - (Required) Template where External-epg to be created.
* `external_epg_name` - (Required) Name of External-epg.
* `contract_name` - (Required) Name of Contract.
* `relationship_type` - (Required) RelationType of the Contract. Values that can be used is provider or consumer
* `contract_schema_id` - (Optional) SchemaID of Contract. schema_id of External-epg will be used if not provided. Should use this parameter when Contract is in different schema than external-epg.
* `contract_template_name` - (Optional) Template Name of Contract. template_name of External-epg will be used if not provided.


## Attribute Reference ##

No attributes are exported.
