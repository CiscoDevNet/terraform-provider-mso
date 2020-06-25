---
layout: "mso"
page_title: "MSO: mso_schema_template_external_epg_subnet"
sidebar_current: "docs-mso-resource-schema_template_external_epg_subnet"
description: |-
  Manages MSO Schema Template External EPG Subnet.
---

# mso_schema_template_external_epg_subnet #

Manages MSO Schema Template External EPG Subnet.

## Example Usage ##

```hcl
resource "mso_schema_template_external_epg_subnet" "subnet1" {
  schema_id = "5ea809672c00003bc40a2799"
  template_name = "Template1"
  external_epg_name =  "UntitledExternalEPG1"
  ip = "10.101.100.0/0"
  name = "sddfgbany"
  scope = ["shared-rtctrl", "export-rtctrl"]
  aggregate = ["shared-rtctrl", "export-rtctrl"]
}
```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy External EPG Subnet.
* `template_name` - (Required) Template where External EPG Subnet to be created.
* `external_epg_name` - (Required) Name of External EPG.
* `ip` - (Required) The IP range in CIDR notation.
* `name` - (Optional) Name of Subnet.
* `scope` - (Optional) The scope of the subnet. Allowed values are `shared-rtctrl`, `export-rtctrl`, `shared-security`, `import-rtctrl`.
* `aggregate` - (Optional) The aggregate of the subnet. Allowed values are `shared-rtctrl`, `export-rtctrl`, `shared-security`, `import-rtctrl`. Aggregate should be enabled only if shared-rtctrl is enabled in Scope.

## Attribute Reference ##

No attributes are exported.
