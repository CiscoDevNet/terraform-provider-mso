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
  schema_id         = mso_schema.schema1.id
  template_name     = "Template1"
  contract_name     = mso_schema_template_contract.template_contract.contract_name
  external_epg_name = mso_schema_template_external_epg.template_externalepg.external_epg_name
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

## Importing ##

An existing MSO Schema Template External Endpoint Group Contract can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_external_epg_contract.c1 {schema_id}/templates/{template_name}/externalEpgs/{external_epg_name}/contractRelationships/{contract_name}/{relationship_type}
```