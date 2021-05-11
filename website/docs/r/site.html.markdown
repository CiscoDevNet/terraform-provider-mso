---
layout: "mso"
page_title: "MSO: mso_site"
sidebar_current: "docs-mso-resource-site"
description: |-
  Manages MSO Site
---

# mso_site #

Manages MSO Site

## Example Usage ##

```hcl
resource "mso_site" "foo_site" {
  name = "mso"
  username = "admin"
  password = "dummypass"
  apic_site_id = "102"
  urls = [ "mso_host" ]
}
```

## Argument Reference ##

* `name` - (Required) The name of the site.
* `username` - (Required) The username for the APICs.
* `password` - (Required) The password for the APICs.
* `apic_site_id` - (Required) The site ID of the APICs.
* `login_domain` - (Optional) Name of login domain. This parameter should be used to authenticate remote user with APIC.
* `maintenance_mode` - (Optional) Boolean flag to enable/disable Maintenance Mode on the site. This parameter is supported only in MSO version 3.0 or higher.
* `urls` - (Required) A list of URLs to reference the APICs.
* `labels` - (Optional) The labels for this site.
* `locations` - (Optional) Location of the site.

## Attribute Reference ##

No Attributes are Exported.

## Importing ##

An existing MSO Site can be [imported][docs-import] into this resource via its Id, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_site.site1 {site_id}
```
