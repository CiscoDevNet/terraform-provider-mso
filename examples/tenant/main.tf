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

resource "mso_tenant" "tenant02" {
  name         = "TangoCh"
  display_name = "Tango"
  description  = "DemoTenant"
  site_associations {
    site_id                = "5c7c95b25100008f01c1ee3c"
    vendor                 = "aws"
    aws_account_id         = "123456789124"
    is_aws_account_trusted = false
    aws_access_key_id      = "AKIAIXCL6LOFME6SUH12"
    aws_secret_key         = "W1fQMYsGKOeK2cJMAnYSpX6uXVP2BrYL8+5uFt23"

  }
  user_associations {
    user_id = "0000ffff0000000000000020"
  }
}

resource "mso_tenant" "tenant01" {
  name         = "TangoCh"
  display_name = "Tango"
  description  = "DemoTenant"
  site_associations {
    site_id                   = "5ce2de773700006a008a2678"
    vendor                    = "azure"
    azure_subscription_id     = "subidtf"
    azure_access_type         = "managed"
    azure_application_id      = "appidtf"
    azure_client_secret       = "clitf"
    azure_active_directory_id = "adidtf"


  }
  user_associations {
    user_id = "0000ffff0000000000000020"
  }
}
