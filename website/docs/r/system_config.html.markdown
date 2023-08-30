---
layout: "mso"
page_title: "MSO: mso_system_config"
sidebar_current: "docs-mso-resource-system_config"
description: |-
  Manages MSO System Configuration.
---

# mso_schema_site_vrf_region #

Manages MSO System Configuration.

## Example Usage ##

```hcl

resource "mso_system_config" "system_config" {
  alias = "test alias"
  banner {
    message = "test message"
    state = "active"
    type = "warning"
  }
  change_control = {
    enable = "enabled"
    number_of_approvers = 2
  }
}

```

## Argument Reference ##

* `alias` - (Optional) The system Alias.
* `banner` - (Optional) A list of Banner configuration. 
    * `state` - (Required) The state of the Banner. Allowed values are `active` or `inactive`.
    * `type` - (Required) The type of the Banner. Allowed values are `critical`, `warning` or `informational`.
    * `message` - (Required) The message of the Banner.
* `change_control` - (Optional) A map of Change Control configuration. 
    * `enable` - (Required) Whether Change Control is enabled. Allowed values are `enabled`, or `disabled`.
    * `number_of_approvers` - (Optional) The number of approvers for the Change Control. MSO defaults to `1` when not provided.

## Attribute Reference ##

No attributes are exported.
