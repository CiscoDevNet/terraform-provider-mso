---
layout: "mso"
page_title: "MSO: mso_l3out_annotation Data Source"
subcategory: "L3Out"
sidebar_current: "docs-mso-data-source-l3out_annotation"
description: |-
  Reads an annotation from an L3Out in NDO.
---

# mso_l3out_annotation (Data Source)

Reads one annotation by key from an L3Out in Nexus Dashboard Orchestration (NDO). The template, L3Out, and annotation must already exist.

- **Supported versions**: Nexus Dashboard 3.1+ (NDO 4.3+).

## API Information

- **Read**: `GET /mso/api/v1/templates/{template_id}` and select the L3Out by UUID and annotation by key.

## GUI Information

In Nexus Dashboard 4.1+ (NDO 5.1+), choose **Manage > Orchestration > Tenant Templates > L3Out**. On earlier supported versions, choose **Configure > Tenant Templates > L3Out**. Open the template and select the L3Out to view its annotations.

## Example Usage

```terraform
data "mso_l3out_annotation" "owner" {
  template_id = mso_template.l3out.id
  l3out_uuid  = mso_l3out.example.uuid
  key         = mso_l3out_annotation.owner.key
}
```

See the [full annotation example](https://github.com/CiscoDevNet/terraform-provider-mso/tree/develop/examples/l3out_annotation) for details on relationships to other Terraform resources and data sources.

## Schema

### Required

- `template_id` (String) The UUID of the L3Out template.
    - **API**: template `templateId`, request path `{template_id}`.
- `l3out_uuid` (String) The UUID of the parent L3Out.
    - **API**: `l3outTemplate.l3outs[].uuid`.
- `key` (String) The annotation key used to select an entry on the parent L3Out.
    - **API**: `l3outTemplate.l3outs[].tagAnnotations[].tagKey`.

### Read-Only

- `value` (String) The annotation value.
    - **API**: `l3outTemplate.l3outs[].tagAnnotations[].tagValue`.
- `id` (String) The Terraform identifier composed of the template UUID, L3Out UUID, and annotation key.

## References

### Related Terraform Objects

- [mso_l3out_annotation resource](/docs/providers/mso/r/l3out_annotation.html)
- [mso_l3out data source](/docs/providers/mso/d/l3out.html)

### External Documentation

- [Cisco Nexus Dashboard 4.3.1 L3Out guide](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-431/external-connectivity-l3out.html)
- [NDO API reference](https://developer.cisco.com/docs/nexus-dashboard/latest/api-reference-orchestration-api-orchestration-api-overview/)
