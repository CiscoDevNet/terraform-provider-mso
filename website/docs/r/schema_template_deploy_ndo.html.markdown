---
layout: "mso"
page_title: "MSO: mso_schema_template_deploy_ndo"
sidebar_current: "docs-mso-resource-schema_template_deploy"
description: |-
  Manages deploy and redeploy operations for schema template.
---

# mso_schema_template_deploy_ndo #

Manages deploy and redeploy operations of schema templates for NDO v3.7 and higher.

## Example Usage ##

```hcl

provider "mso" {
  username = "" # <MSO username>
  password = "" # <MSO pwd>
  url      = "" # <MSO URL>
  insecure = true
  platform = "nd"
}

resource "mso_schema_template_deploy_ndo" "template_deployer" {
  schema_id     = mso_schema.schema1.id
  template_name = "Template1"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema-id of the template.
* `template_name` - (Required) The name of the template to deploy or redeploy.
* `re_deploy` - (Optional) Boolean flag indicating whether to re-deploy the template to the associated sites. Default is false, which would trigger a regular deploy operation. 

### Notes ###

* This resource requires 'platform = "nd"' to be configured in the provider configuration section.
* This resource is intentionally created non-idempotent so that it deploys the template in every run, it will not fail if there is no change and we deploy or redeploy the template again. When destroying the resource, no action is taken.
* Prior to deploy or redeploy a schema validation is executed. When schema validation fails, the resource will fail and deploy or redeploy will not be executed.
* A template can only be undeployed from a site by disassociating the site from the template with the resource mso_schema_site.

## Attribute Reference ##

No attributes are exported.
