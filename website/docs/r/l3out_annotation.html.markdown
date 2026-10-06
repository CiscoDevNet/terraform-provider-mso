---
layout: "mso"
page_title: "MSO: mso_l3out_annotation"
subcategory: "L3Out"
sidebar_current: "docs-mso-resource-l3out_annotation"
description: |-
  Manages one annotation on an L3Out in Nexus Dashboard Orchestration.
---

# mso_l3out_annotation (Resource)

Manages one annotation identified by its key on an existing L3Out. The template and L3Out must exist. Separate resources can manage different keys on the same L3Out.

-> Leave `mso_l3out.annotations` omitted when managing an annotation with this resource. A configured parent map owns the entire annotation collection.

## API Information

- **APIs**: Nexus Dashboard Orchestration API (template endpoints).
- **Create**: `GET /mso/api/v1/templates/{template_id}`, then `PATCH /mso/api/v1/templates/{template_id}` to add a `tagAnnotations` entry.
- **Read**: `GET /mso/api/v1/templates/{template_id}` and select the L3Out by UUID and annotation by key.
- **Update**: `GET` the current template, then `PATCH /mso/api/v1/templates/{template_id}` to change `tagValue`.
- **Delete**: `GET` the current template, then `PATCH /mso/api/v1/templates/{template_id}` to remove the selected entry.
- **Target versions**: Nexus Dashboard 4.1+ (Orchestration 5.1+).
- **Terraform ID**: `{template_id}/{l3out_uuid}/{key}`.
- **Managed resource identity**: `template_id`, `l3out_uuid`, and `key`.

The provider sends template PATCH requests with `validate=false` and refreshes the selected object from the template GET response.

## GUI Information

In Nexus Dashboard, choose **Manage > Orchestration > Tenant Templates > L3Out**, then open the template and select the parent L3Out. The provider manages that object's `tagAnnotations` entries through the API. Cisco documents the [L3Out template workflow](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-411/external-connectivity-l3out.html).

## Example Usage

```terraform
resource "mso_l3out_annotation" "owner" {
  template_id = mso_template.l3out.id
  l3out_uuid  = mso_l3out.example.uuid
  key         = "owner"
  value       = "network"
}
```

The [complete standalone annotation example](https://github.com/CiscoDevNet/terraform-provider-mso/tree/master/examples/l3out_annotation) creates the parent L3Out, manages two keys, and reads one key through the data source.

## Schema

### Required

- `template_id` (String, API: template `templateId`, request path `{template_id}`) The UUID of the existing L3Out template.
    - **Force New**: `true`.
- `l3out_uuid` (String, API: `l3outTemplate.l3outs[].uuid`) The UUID of the parent L3Out.
    - **Force New**: `true`.
- `key` (String, API: `l3outTemplate.l3outs[].tagAnnotations[].tagKey`) The unique annotation key on the parent L3Out.
    - **Validation**: Must contain at least one character. Orchestration permits only one annotation per key.
    - **Force New**: `true`; changing the key creates a different annotation.
- `value` (String, API: `l3outTemplate.l3outs[].tagAnnotations[].tagValue`) The value of the annotation. Changing it updates the selected annotation.
    - **Validation**: An empty string is allowed. Different keys may have the same value.

### Read-Only

- `id` (String, provider generated) The stable Terraform ID composed of the template UUID, L3Out UUID, and annotation key.

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

- [Cisco Nexus Dashboard 4.1.1 L3Out guide](https://www.cisco.com/c/en/us/td/docs/dcn/nd/4x/articles-411/external-connectivity-l3out.html)
- [Cisco Nexus Dashboard Orchestration Get Template API](https://developer.cisco.com/docs/nexus-dashboard/latest/get-template/)
- [Cisco Nexus Dashboard Orchestration Patch Template API](https://developer.cisco.com/docs/nexus-dashboard/latest/patch-template/)
- [Cisco Nexus Dashboard 4.1.1 API introduction](https://developer.cisco.com/docs/nexus-dashboard/4-1/)
