---
layout: "mso"
page_title: "MSO: mso_site"
sidebar_current: "docs-mso-data-source-site"
description: |-
  Data source for MSO Site
---

# mso_site #

Data source for MSO site

## Example Usage ##

```hcl
data "mso_site" "sample_site" {
  name  = "AWS-West"
}
```

## Argument Reference ##

* `name` - (Required) The name of the site.

## Attribute Reference ##

* `username` - (Optional) The username for the APICs.
* `password` - (Optional) The password for the APICs.
* `apic_site_id` - (Optional) The site ID of the APICs.
* `urls` - (Optional) A list of URLs to reference the APICs.
* `labels` - (Optional) The labels for this site.
* `location` - (Optional) Location of the site.
