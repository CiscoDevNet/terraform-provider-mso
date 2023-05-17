---
layout: "mso"
page_title: "MSO: mso_tenant"
sidebar_current: "docs-mso-resource-tenant"
description: |-
  Manages MSO Tenant
---

# mso_tenant #

Manages MSO Tenant

## Example Usage ##

```hcl

data "mso_site" "site1" {
  name = "site1"
}

data "mso_user" "user1" {
  username = "user1"
}

# With No Site Association
resource "mso_tenant" "tenant1" {
  name         = "tenant1"
  display_name = "tenant1"
  site_associations { 
    site_id = mso_site.site1.id 
  }
  user_associations { 
    user_id = mso_user.user1.id 
  }
}

# With AWS Site Association
resource "mso_tenant" "tenant2" {
  name         = "tenant2"
  display_name = "tenant2"
  description  = "demo tenant 2"
  site_associations {
    site_id                = mso_site.site1.id
    vendor                 = "aws"
    aws_account_id         = "123456789124"
    is_aws_account_trusted = false
    aws_access_key_id      = "AKIAIXCL6LOFME6SUH12"
    aws_secret_key         = "W1fQMYsGKOeK2cJMAnYSpX6uXVP2BrYL8+5uFt23"
  }
  user_associations {
    user_id = mso_user.user1.id
  }
}

# With Azure Site Association
resource "mso_tenant" "tenant3" {
  name         = "tenant3"
  display_name = "tenant3"
  description  = "demo tenant 3"
  site_associations {
    site_id                   = mso_site.site1.id
    vendor                    = "azure"
    azure_subscription_id     = "subidtf"
    azure_access_type         = "credentials"
    azure_application_id      = "appidtf"
    azure_client_secret       = "clitf"
    azure_active_directory_id = "adidtf"
  }
  user_associations {
    user_id = mso_user.user1.id
  }
}

# With GCP Site Association
resource "mso_tenant" "tenant4" {
  name         = "tenant4"
  display_name = "tenant4"
  description  = "demo tenant 4"
  site_associations {
    site_id         = data.mso_site.demo_site.id
    vendor          = "gcp"
    gcp_project_id  = "10"
    gcp_access_type = "unmanaged"
    gcp_email       = "demo@tenant.com"
		gcp_name        = "demo_name"
		gcp_key_id      = "demo_key"
    gcp_private_key = "demo_private_key"
    gcp_client_id   = "demo_client_id"
  }
  user_associations {
    user_id = mso_user.user1.id
  }
}

```

## Argument Reference ##

* `name` - (Required) The name of the tenant.
* `display_name` - (Required) The name of the tenant to be displayed in the web UI.
* `description` - (Optional) The description for this tenant.
* `orchestrator_only` - (Optional) Option to delete this tenant only from orchestrator or not. Default value is "false".
* `user_associations` - (Optional) A list of associated users for this tenant.
* `user_associations.user_id` - (Optional) Id of user to be associated to this tenant.
* `site_association` - (Optional) A list of associated sites for this tenant.
* `site_association.site_id` - (Optional) Id of site to associate with this Tenant.
* `site_association.security_domains` - (Optional) Security domains to associate with this Site.
* `site_association.vendor` - (Optional) Name of cloud vendor in the case of Attaching cloud site. Allowed values are `aws`, `azure`, and `gcp`.
* `site_association.aws_account_id` - (Optional) Id of AWS account. It's required when vendor is set to aws. This parameter will only have effect with `vendor` = aws
* `site_association.is_aws_account_trusted` - (Optional) Boolean flag to indicate whether this account is trusted or not. Trusted account does not require aws_access_key_id and aws_secret_key.
* `site_association.aws_access_key_id` - (Optional) AWS Access Key Id. It must be provided if the AWS account is not trusted. This parameter will only have effect with `vendor` = aws.
* `site_association.aws_secret_key` - (Optional) AWS Secret Key Id. It must be provided if the AWS account is not trusted. This parameter will only have effect with `vendor` = aws.
* `site_association.azure_subscription_id` - (Optional) Azure subscription id. It's required when vendor is set to azure. This parameter will only have effect with `vendor` = azure.
* `site_association.azure_access_type` - (Optional) Type of Azure Account Configuration. Allowed values are `managed`, `shared` and `credentials`. Default to `managed`. Other Credentials are not required if azure_access_type is set to managed. This parameter will only have effect with `vendor` = azure.
* `site_association.azure_application_id` - (Optional) Azure Application Id. It must be provided when azure_access_type to credentials. This parameter will only have effect with `vendor` = azure.
* `site_association.azure_client_secret` - (Optional) Azure Client Secret. It must be provided when azure_access_type to credentials. This parameter will only have effect with `vendor` = azure.
* `site_association.azure_active_directory_id` - (Optional) Azure Active Directory Id. It must be provided when azure_access_type to credentials. This parameter will only have effect with `vendor` = azure.
* `site_association.azure_shared_account_id` - (Optional) Azure shared account Id. It must be provided when azure_access_type to shared. This parameter will only have effect with `vendor` = azure.
* `site_association.gcp_project_id` - (Optional) GCP Project Id. It must be provided for the GCP account. This parameter will only have effect with `vendor` = gcp.
* `site_association.gcp_access_type` - (Optional) Type of GCP Account Configuration. Allowed values are `managed` or `unmanaged`. This parameter will only have effect with `vendor` = gcp.
* `site_association.gcp_name` - (Optional) GCP Name. It must be provided if the GCP account is not managed. This parameter will only have effect with `vendor` = gcp.
* `site_association.gcp_key_id` - (Optional) GCP Key Id. It must be provided if the GCP account is not managed. This parameter will only have effect with `vendor` = gcp.
* `site_association.gcp_private_key` - (Optional) GCP Private Key. This parameter will only have effect with `vendor` = gcp.
* `site_association.gcp_client_id` - (Optional) GCP Client Id. This parameter will only have effect with `vendor` = gcp.
* `site_association.gcp_email` - (Optional) GCP Email. It must be provided if the GCP account is not managed. This parameter will only have effect with `vendor` = gcp.

NOTE: AWS, Azure or GCP credentials will be used based on whatever is passed in `vendor` argument if both (AWS + Azure + GCP) Credentials are provided.

## Attribute Reference ##

The only Attribute exposed for this resource is `id`. Which is set to the id of tenant created.

## Importing ##

An existing MSO Tenant can be [imported][docs-import] into this resource via its Id, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_tenant.tenant1 {tenant_id}
```