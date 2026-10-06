---
layout: "mso"
page_title: "MSO: mso_l3out_annotation Data Source"
subcategory: "L3Out"
sidebar_current: "docs-mso-data-source-l3out_annotation"
description: |-
  Reads one annotation from an L3Out in Nexus Dashboard Orchestration.
---

# mso_l3out_annotation (Data Source)

Reads one annotation by key from an L3Out in a Nexus Dashboard Orchestration L3Out template. The template, L3Out, and annotation must already exist.

## API Information

- **APIs**: Nexus Dashboard Orchestration API (template endpoints).
- **Read**: `GET /mso/api/v1/templates/{template_id}`; select the L3Out by UUID and annotation by `tagKey`.
- **Target versions**: Nexus Dashboard 4.1+ (Orchestration 5.1+).
- **Terraform ID**: `{template_id}/{l3out_uuid}/{key}`.

## GUI Information

In Nexus Dashboard, choose **Manage > Orchestration > Tenant Templates > L3Out**, then open the template and select the parent L3Out. This data source reads its `tagAnnotations` entries through the API. See Cisco's [L3Out template workflow](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-411/external-connectivity-l3out.html).

## Example Usage

```terraform
data "mso_l3out_annotation" "owner" {
  template_id = mso_template.l3out.id
  l3out_uuid  = mso_l3out.example.uuid
  key         = mso_l3out_annotation.owner.key
}
```

The [complete standalone annotation example](https://github.com/CiscoDevNet/terraform-provider-mso/tree/master/examples/l3out_annotation) includes the parent resources and this lookup.

## Schema

### Required

- `template_id` (String, API: template `templateId`, request path `{template_id}`) The UUID of the L3Out template to read.
- `l3out_uuid` (String, API: `l3outTemplate.l3outs[].uuid`) The UUID of the parent L3Out.
- `key` (String, API: `l3outTemplate.l3outs[].tagAnnotations[].tagKey`) The annotation key used to select an entry on the parent L3Out.
    - **Validation**: Must contain at least one character.

### Optional

None.

### Read-Only

- `value` (String, API: `l3outTemplate.l3outs[].tagAnnotations[].tagValue`) The value of the annotation, including an empty string when configured that way.
- `id` (String, provider generated) The stable identifier composed of the template UUID, L3Out UUID, and annotation key.

## References

### Related Terraform Objects

- [mso_l3out_annotation resource](/docs/providers/mso/r/l3out_annotation.html)
- [mso_l3out data source](/docs/providers/mso/d/l3out.html)

### External Documentation

- [Cisco Nexus Dashboard 4.1.1 L3Out guide](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-411/external-connectivity-l3out.html)
- [Cisco Nexus Dashboard Orchestration Get Template API](https://developer.cisco.com/docs/nexus-dashboard/latest/get-template/)
- [Cisco Nexus Dashboard 4.1.1 API introduction](https://developer.cisco.com/docs/nexus-dashboard/4-1/)
