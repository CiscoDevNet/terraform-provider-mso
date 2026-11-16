---
layout: "mso"
page_title: "MSO: mso_tenant"
sidebar_current: "docs-mso-data-source-tenant"
description: |-
  Data source for MSO Tenant.
---

# mso_tenant #

Data source for MSO Tenant.

## Example Usage ##

```hcl

data "mso_tenant" "example" {
  name = "mso"
}

```

## Argument Reference ##

* `name` - (Required) The name of the Tenant.

## Attribute Reference ##

* `display_name` - (Read-Only) The name of the Tenant as displayed on the MSO UI.
* `description` - (Read-Only) The description of the Tenant.
* `user_associations` - (Read-Only) A list of associated users of the Tenant.
    * `user_id` - (Read-Only) The user ID associated to this tenant.
* `site_association` - (Read-Only) A list of associated sites of the Tenant.
    * `site_id` - (Read-Only) The site ID associated with this Tenant.
    * `security_domains` - (Read-Only) The security domain associated with this Tenant.
    * `vendor` - (Read-Only) The cloud vendor associated with this Tenant. Only applicable for cloud sites.
    * `aws_account_id` - (Read-Only) The ID of the AWS account.
    * `is_aws_account_trusted` - (Read-Only) Whether this account is trusted.
    * `aws_access_key_id` - (Read-Only) The Access Key ID of the AWS account.
    * `aws_secret_key` - (Read-Only) The Secret Key ID of the AWS account.
    * `azure_subscription_id` - (Read-Only) The subscription ID of the Azure account.
    * `azure_access_type` - (Read-Only) The type of the Azure account.
    * `azure_application_id` - (Read-Only) The application ID of the Azure account.
    * `azure_client_secret` - (Read-Only) The client secret of the Azure account.
    * `azure_active_directory_id` - (Read-Only) The active directory ID of the Azure account.
    * `azure_shared_account_id` - (Read-Only) The shared account ID of the Azure account.
    * `gcp_project_id` - (Read-Only) The project ID of the GCP account.
    * `gcp_access_type` - (Read-Only) The access type of the GCP account.
    * `gcp_name` - (Read-Only) The name of the GCP account.
    * `gcp_key_id` - (Read-Only) The key ID of the GCP account.
    * `gcp_private_key` - (Read-Only) The private key of the GCP account.
    * `gcp_client_id` - (Read-Only) The client ID of the GCP account.
    * `gcp_email` - (Read-Only) The email of the GCP account.
