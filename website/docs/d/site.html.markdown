---
layout: "mso"
page_title: "MSO: mso_site"
sidebar_current: "docs-mso-data-source-site"
description: |-
  Data source for MSO Site.
---

# mso_site #

Data source for MSO Site.

## Example Usage ##

```hcl

data "mso_site" "example" {
  name = "AWS-West"
}

```

## Argument Reference ##

* `name` - (Required) The name of the Site.

## Attribute Reference ##

* `username` - (Read-Only) The username of the Site.
* `password` - (Read-Only) The password of the Site.
* `type` - (Read-Only) The type of the Site.
* `group_id` - (Read-Only) The group ID of the Site.
* `version` - (Read-Only) The software version of the Site.
* `status` - (Read-Only) The connectivity status of the Site.
* `reprovision` - (Read-Only) Whether the Site needs a reprovision.
* `proxy` - (Read-Only) Whether the Site uses a proxy.
* `sr_l3out` - (Read-Only) Whether the Site has segment routing l3out enabled.
* `template_count` - (Read-Only) The amount of templates attached to the Site.
* `apic_site_id` - (Read-Only) The ID of the Site.
* `cloud_providers` - (Read-Only) A list of cloud providers for the Site.
* `urls` - (Read-Only) A list of URLs to reference the Site.
* `labels` - (Read-Only) The labels of the Site.
* `location` - (Read-Only) The location of the Site.
    * `lat` - (Read-Only) The latitude of the Site.
    * `long` - (Read-Only) The longitude of the Site.
