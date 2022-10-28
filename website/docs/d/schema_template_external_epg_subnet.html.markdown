---
layout: "mso"
page_title: "MSO: mso_schema_template_external_epg_subnet"
sidebar_current: "docs-mso-data-source-schema_template_external_epg_subnet"
description: |-
  Data source for MSO Schema Template External EPG Subnet.
---

# mso_schema_template_external_epg_subnet #

Data source for MSO Schema Template External EPG Subnet.

## Example Usage ##

```hcl

data "mso_schema_template_external_epg_subnet" "subnet1" {
  schema_id         = data.mso_schema.schema1.id
  template_name     = "Template1"
  external_epg_name = "Internet"
  ip                = "30.1.1.0/24"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy External EPG Subnet.
* `template_name` - (Required) Template where External EPG Subnet to be created.
* `external_epg_name` - (Required) Name of External EPG.
* `ip` - (Required) The IP range in CIDR notation.

## Attribute Reference ##

* `name` - (Optional) Name of Subnet.
* `scope` - (Optional) The scope of the subnet. Allowed values are `shared-rtctrl`, `export-rtctrl`, `shared-security`, `import-rtctrl`, `import-security`.
* `aggregate` - (Optional) The aggregate of the subnet. Allowed values are `shared-rtctrl`, `export-rtctrl`, `shared-security`, `import-rtctrl`. Aggregate should be enabled only if shared-rtctrl is enabled in Scope.
