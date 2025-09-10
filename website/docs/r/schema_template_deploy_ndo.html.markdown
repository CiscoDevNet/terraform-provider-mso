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
  retries  = 3
}

# Deploying an Application Template by Name
resource "mso_schema_template_deploy_ndo" "template_deployer" {
  schema_id     = mso_schema.schema1.id
  template_name = "Template1"
}

# Deploying any Template by ID
resource "mso_schema_template_deploy_ndo" "deploy_by_id" {
  template_id = "68b616a4d3bd0f48316c176b"
}

# Deploying a Template of type tenant policy
resource "mso_schema_template_deploy_ndo" "deploy_by_id" {
  template_name = "Template1"
  template_type = "tenant"
}
```

## Argument Reference ##

* `schema_id` - (Optional) The ID of the schema that contains the template. This is required when deploying an application-type template by name.
* `template_name` - (Optional) The name of the template to deploy. This is required when identifying a template by name instead of by its template_id.
* `template_id` - (Optional) The unique ID of the template to deploy. If this is provided, it takes precedence over schema_id and template_name.
* `template_type` - (Optional) The type of the template. This is used in combination with template_name to uniquely identify a non-application template. Default is application.
* `re_deploy` - (Optional) Boolean flag indicating whether to re-deploy the template to the associated sites. Default is false, which would trigger a regular deploy operation. 

### Notes ###

* This resource requires 'platform = "nd"' to be configured in the provider configuration section.
* This resource is intentionally created non-idempotent so that it deploys the template in every run, it will not fail if there is no change and we deploy or redeploy the template again. When destroying the resource, no action is taken.
* Prior to deploy or redeploy a schema validation is executed. When schema validation fails, the resource will fail and deploy or redeploy will not be executed.
* A template can only be undeployed from a site by disassociating the site from the template with the resource mso_schema_site.
* To adjust the number of retries to ensure successful deployment completion, configure the retries argument in the provider configuration section.

## Attribute Reference ##

No attributes are exported.
