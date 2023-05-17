terraform {
  required_providers {
    mso = {
      source = "CiscoDevNet/mso"
    }
  }
}

provider "mso" {
  username = "" # <MSO username>
  password = "" # <MSO pwd>
  url      = "" # <MSO URL>
  insecure = true
}

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
resource "mso_tenant" "aws_tenant" {
  name         = "aws_tenant"
  display_name = "aws_tenant"
  description  = "demo tenant with aws site association"
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
resource "mso_tenant" "azure_tenant" {
  name         = "azure_tenant"
  display_name = "azure_tenant"
  description  = "demo tenant with azure site association"
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
resource "mso_tenant" "gcp_tenant" {
  name         = "gcp_tenant"
  display_name = "gcp_tenant"
  description  = "demo tenant with gcp site association"
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
