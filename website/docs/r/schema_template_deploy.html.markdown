---
layout: "mso"
page_title: "MSO: mso_schema_template_deploy"
sidebar_current: "docs-mso-resource-schema_template_deploy"
description: |-
  Manages deploy/undeploy operations for schema template on sites.
---

# mso_schema_template_deploy #

Manages deploy/undeploy operations for schema template on sites.

## Example Usage ##

```hcl

resource "mso_schema_template_deploy" "template_deployer" {
  schema_id     = mso_schema.schema1.id
  template_name = "Template1"
  site_id       = mso_site.site1.id
  undeploy      = true
}

```

## Argument Reference ##


* `schema_id` - (Required) The schema-id of template.
* `template_name` - (Required) name of the template to deploy/undeploy.
* `undeploy` - (Optional) Boolean flag indicating whether to undeploy the template from a single site (see site_id) or not. Default is false.
* `site_id` - (Optional) Site-id from where the template is to be undeployed. It is required if you set undeploy = true.

NOTE: This resource is intentionally created non-idempotent so that it deploys the template in every run, it will not fail if there is no change and we deploy the template again. When destroying the resource, all sites will be undeployed.


## Attribute Reference ##

No attributes are exported.



