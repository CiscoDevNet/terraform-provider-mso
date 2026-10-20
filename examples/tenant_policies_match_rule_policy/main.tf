terraform {
  required_providers {
    mso = {
      source = "CiscoDevNet/mso"
    }
  }
}

provider "mso" {
  username = "" # <MSO username>
  password = "" # <MSO password>
  url      = "" # <MSO URL>
  insecure = true
}

data "mso_tenant" "example" {
  name = "example_tenant"
}

resource "mso_template" "tenant" {
  template_name = "tenant_template"
  template_type = "tenant"
  tenant_id     = data.mso_tenant.example.id
}

resource "mso_tenant_policies_match_rule_policy" "example" {
  template_id = mso_template.tenant.id
  name        = "example_match_rule_policy"
  description = "Match selected routes"

  match_prefix {
    prefix      = "10.10.10.0/24"
    description = "IPv4 routes"
    aggregate   = true
    from_prefix = 25
    to_prefix   = 32
  }

  match_prefix {
    prefix      = "2001:db8::/64"
    description = "IPv6 routes"
    aggregate   = false
    from_prefix = 0
    to_prefix   = 0
  }

  match_community_term {
    name        = "export_community"
    description = "Match export communities"

    match_community_factor {
      community   = "no-export"
      scope       = "transitive"
      description = "Match no-export"
    }

    match_community_factor {
      community   = "no-advertise"
      scope       = "non-transitive"
      description = "Match no-advertise"
    }
  }
}

data "mso_tenant_policies_match_rule_policy" "example" {
  template_id = mso_tenant_policies_match_rule_policy.example.template_id
  name        = mso_tenant_policies_match_rule_policy.example.name
}
