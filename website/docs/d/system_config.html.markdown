---
layout: "mso"
page_title: "MSO: mso_system_config"
sidebar_current: "docs-mso-data-source-system_config"
description: |-
  Data source for MSO System Configuration.
---

# mso_schema_site_vrf_region #

Data source for MSO System Configuration.

## Example Usage ##

```hcl

data "mso_system_config" "system_config" {}

```

## Argument Reference ##

No arguments are required.

## Attribute Reference ##

* `alias` - (Read-Only) The system Alias.
* `banner` - (Read-Only) A list of Banner configuration.
    * `state` - (Read-Only) The state of the Banner.
    * `type` - (Read-Only) The type of the Banner.
    * `message` - (Read-Only) The message of the Banner.
* `change_control` - (Read-Only) A map of Change Control configuration.
    * `workflow` - (Required) Whether Change Control workflow is enabled. 
    * `number_of_approvers` - (Read-Only) The number of approvers for the Change Control. 
