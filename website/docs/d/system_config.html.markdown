---
layout: "mso"
page_title: "MSO: mso_system_config"
sidebar_current: "docs-mso-data-source-system_config"
description: |-
  Data source for MSO System Configuration.
---

# mso_system_config #

!> **Deprecated** This data source is deprecated: no longer functional on Nexus Dashboard (ND) 4.0+ / NDO 5.0+ and will be removed once ND 3.x / NDO 4.x is no longer supported.

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
    * `workflow` - (Read-Only) Whether Change Control workflow is enabled. 
    * `number_of_approvers` - (Read-Only) The number of approvers for the Change Control. 
