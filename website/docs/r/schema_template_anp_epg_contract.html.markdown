---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg_contract"
sidebar_current: "docs-mso-resource-schema_template_anp_epg_contract"
description: |-
  Manages MSO Schema Template Application Network Profile(ANP) Endpoint Group(EPG) Contract.
---

# mso_schema_template_anp_epg_contract #

Manages MSO Schema Template Application Network Profile Endpoint Group Contract resource.

## Example Usage ##

```hcl
resource "mso_schema_template_anp_epg_contract" "contract1" {
  schema_id = "5c6c16d7270000c710f8094d"
  template_name = "Template1"
  anp_name = "WoS-Cloud-Only-2"
  epg_name = "DB"
  contract_name = "Internet-access"
  relationship_type = "provider"
  
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy ANP EPG Contract.
* `template_name` - (Required) Template name under which you want to deploy ANP EPG Contract.
* `anp_name` - (Required) ANP name under which you want to deploy ANP EPG Contract.
* `epg_name` - (Required) EPG name under which you want to deploy ANP EPG Contract.
* `contract_name` - (Required) The contract name which you want to associate with.
* `relationship_type` - (Required) The type of the contract i.e. provider or consumer.
* `contract_schema_id` - (Optional) SchemaID of Contract. schema_id of ANP EPG will be used if not provided. Should use this parameter when Contract is in different schema than ANP EPG.
* `contract_template_name` - (Optional) Template Name associated with Contract. template_name of ANP EPG will be used if not provided. Should use this parameter when Contract is in different schema than ANP EPG.

## Attribute Reference ##

No attributes are exported.
