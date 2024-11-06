---
layout: "mso"
page_title: "MSO: mso_template"
sidebar_current: "docs-mso-data-source-template"
description: |-
  Data source for MSO Template.
---

# mso_template #

Data source for MSO Template.

## Example Usage ##

```hcl

data "mso_template" "example_with_name" {
  template_name = "tenant_template"
  template_type = "tenant"
}

data "mso_template" "example_with_id" {
  template_id = "6718b46395400f3759523378"
}

```

## Argument Reference ##

* `template_id` - (Optional) The ID of the template. Mutually exclusive with `template_name`.
* `template_name` - (Optional) The name of the template. Mutually exclusive with `template_id`.
* `template_type` - (Optional) The type of the template. Allowed values are `tenant`, `l3out`, `fabric_policy`, `fabric_resource`, `monitoring_tenant`, `monitoring_access`, or `service_device`. Required when `template_name` is provided.

## Attribute Reference ##

* `tenant_id` - (Read-Only) The ID of the tenant associated with the template.
* `sites` - (Read-Only) A list of site names associated with the template.
