---
layout: "mso"
page_title: "MSO: mso_site"
sidebar_current: "docs-mso-resource-site"
description: |-
  Manages MSO Site
---

# schema #

Manages MSO Site

## Example Usage ##

```hcl
resource "mso_site" "foo_site" {
  name = "mso"
  username = "admin"
  password = "noir0!234"
  apic_site_id = "102"
  urls = [ "https://3.208.123.222/" ]
}
```

## Argument Reference ##

* `name` - (Required) The name of the site.
* `username` - (Required) The username for the APICs.
* `password` - (Required) The password for the APICs.
* `apic_site_id` - (Required) The site ID of the APICs.
* `urls` - (Required) A list of URLs to reference the APICs.

## Attribute Reference ##

* `labels` - (Optional) The labels for this site.
* `locations` - (Optional) Location of the site.
