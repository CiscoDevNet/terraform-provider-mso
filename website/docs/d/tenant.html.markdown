---
layout: "mso"
page_title: "MSO: mso_tenant"
sidebar_current: "docs-mso-data-source-tenant"
description: |-
  Data source for MSO Tenant
---

# mso_tenant #

Data source for MSO tenant

## Example Usage ##

```hcl

data "mso_tenant" "tenant1" {
  name = "mso"
}

```

## Argument Reference ##

* `name` - (Required) The name of the tenant.

## Attribute Reference ##

* `display_name` - (Read-Only) The name of the tenant to be displayed in the web UI.
* `description` - (Read-Only) The description for this tenant.
* `orchestrator_only` - (Read-Only) Option to delete this tenant only from orchestrator or not. Default value is "false".
* `user_associations` - (Read-Only) A list of associated users for this tenant.
* `user_associations.user_id` - (Read-Only) Id of user to be associated to this tenant.
* `site_association` - (Read-Only) A list of associated sites for this tenant.
* `site_association.site_id` - (Read-Only) Id of site to associate with this Tenant.
* `site_association.security_domains` - (Read-Only) Security domains to associate with this Site.
* `site_association.vendor` - (Read-Only) Name of cloud vendor in the case of Attaching cloud site.
* `site_association.aws_account_id` - (Read-Only) Id of AWS account.
* `site_association.is_aws_account_trusted` - (Read-Only) Boolean flag to indicate whether this account is trusted or not.
* `site_association.aws_access_key_id` - (Read-Only) AWS Access Key Id.
* `site_association.aws_secret_key` - (Read-Only) AWS Secret Key Id.
* `site_association.azure_subscription_id` - (Read-Only) Azure subscription id.
* `site_association.azure_access_type` - (Read-Only) Type of Azure Account Configuration.
* `site_association.azure_application_id` - (Read-Only) Azure Application Id.
* `site_association.azure_client_secret` - (Read-Only) Azure Client Secret.
* `site_association.azure_active_directory_id` - (Read-Only) Azure Active Directory Id.
* `site_association.azure_shared_account_id` - (Read-Only) Azure shared account Id.
* `site_association.gcp_project_id` - (Read-Only) GCP Project Id. It must be provided for the GCP account.
* `site_association.gcp_access_type` - (Read-Only) Type of GCP Account Configuration.
* `site_association.gcp_name` - (Read-Only) GCP Name.
* `site_association.gcp_key_id` - (Read-Only) GCP Key Id.
* `site_association.gcp_private_key` - (Read-Only) GCP Private Key.
* `site_association.gcp_client_id` - (Read-Only) GCP Client Id.
* `site_association.gcp_email` - (Read-Only) GCP Email.
