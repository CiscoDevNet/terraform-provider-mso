---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg_subnet"
sidebar_current: "docs-mso-resource-schema_template_anp_epg_subnet"
description: |-
  Manages MSO Schema Template Application Network Profiles Endpoint Groups Subnets.
---

# mso_schema_template_anp_epg_subnet#

Manages MSO Schema Template Application Network Profiles Endpoint Groups Subnets.

## Example Usage ##

```hcl
resource "mso_schema_template_anp_epg_subnet" "subnet1" {
  schema_id = "5c6c16d7270000c710f8094d"
  anp_name = "WoS-Cloud-Only-2"
  epg_name ="EPG4"
  template = "Template1"
  ip = "31.101.102.0/8"
  scope = "public"
  shared = true
  
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg Subnet.
* `template_name` - (Required) Template where Anp Epg Subnet to be created.
* `anp_name` - (Required) Name of Application Network Profiles.
* `epg_name` - (Required) Name of Endpoint Group.
* `ip` - (Required) Ip Address of Subnet.
* `scope` - (Optional) Scope of Subnet.
* `shared` - (Optional) Whether the subnet should be shared or not.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Schema Template Application Network Profiles Endpoint Groups Subnet can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_anp_epg_subnet.subnet1 {schema_id}/template/{template_name}/anp/{anp_name}/epg/{epg_name}/subnet/{ip}
```
