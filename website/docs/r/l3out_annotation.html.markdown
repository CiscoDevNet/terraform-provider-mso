---
layout: "mso"
page_title: "MSO: mso_l3out_annotation"
subcategory: "L3Out"
sidebar_current: "docs-mso-resource-l3out_annotation"
description: |-
  Manages an annotation on an L3Out in NDO.
---

# mso_l3out_annotation (Resource)

Manages one annotation by key on an existing L3Out in Nexus Dashboard Orchestration (NDO). Separate resources can manage different annotation keys on the same L3Out.

- **Supported versions**: Nexus Dashboard 3.1+ (NDO 4.3+).

-> Leave `mso_l3out.annotations` omitted when managing annotations with this resource. A configured parent map manages the entire annotation collection.

## API Information

- **Create**: `GET /mso/api/v1/templates/{template_id}`, then `PATCH /mso/api/v1/templates/{template_id}` to add the annotation.
- **Read**: `GET /mso/api/v1/templates/{template_id}` and select the L3Out by UUID and annotation by key.
- **Update**: `GET` the current template, then `PATCH /mso/api/v1/templates/{template_id}` to update the annotation value.
- **Delete**: `GET` the current template, then `PATCH /mso/api/v1/templates/{template_id}` to remove the annotation.

## GUI Information

In Nexus Dashboard 4.1+ (NDO 5.1+), choose **Manage > Orchestration > Tenant Templates > L3Out**. On earlier supported versions, choose **Configure > Tenant Templates > L3Out**. Open the template and select the L3Out to view or edit its annotations.

## Example Usage

```terraform
resource "mso_l3out_annotation" "owner" {
  template_id = mso_template.l3out.id
  l3out_uuid  = mso_l3out.example.uuid
  key         = "owner"
  value       = "network"
}
```

See the [full annotation example](https://github.com/CiscoDevNet/terraform-provider-mso/tree/develop/examples/l3out_annotation) for details on relationships to other Terraform resources and data sources.

## Schema

### Required

- `template_id` (String) The UUID of the existing L3Out template.
    - **API**: template `templateId`, request path `{template_id}`.
    - **Force New**: `true`.
- `l3out_uuid` (String) The UUID of the parent L3Out.
    - **API**: `l3outTemplate.l3outs[].uuid`.
    - **Force New**: `true`.
- `key` (String) The annotation key, unique within the parent L3Out.
    - **API**: `l3outTemplate.l3outs[].tagAnnotations[].tagKey`.
    - **Validation**: Must contain at least one character.
    - **Force New**: `true`.
- `value` (String) The annotation value. Different keys can have the same value.
    - **API**: `l3outTemplate.l3outs[].tagAnnotations[].tagValue`.
    - **Validation**: NDO requires a nonempty annotation value.

### Read-Only

- `id` (String) The Terraform identifier composed of the template UUID, L3Out UUID, and annotation key.

## Importing

```shell
terraform import mso_l3out_annotation.owner '<template_id>/<l3out_uuid>/<key>'
```

Terraform 1.12 and later can also import with the managed resource identity:

```terraform
import {
  to = mso_l3out_annotation.owner
  identity = {
    template_id = "<template_id>"
    l3out_uuid  = "<l3out_uuid>"
    key         = "owner"
  }
}
```

## References

### Related Terraform Objects

- [mso_l3out resource](/docs/providers/mso/r/l3out.html)
- [mso_l3out_annotation data source](/docs/providers/mso/d/l3out_annotation.html)

### External Documentation

- [Cisco Nexus Dashboard 4.3.1 L3Out guide](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-431/external-connectivity-l3out.html)
- [NDO API reference](https://developer.cisco.com/docs/nexus-dashboard/latest/api-reference-orchestration-api-orchestration-api-overview/)
