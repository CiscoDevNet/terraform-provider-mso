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

resource "mso_tenant_policies_set_rule_policy" "example" {
  template_id = mso_template.tenant.id
  name        = "example_set_rule_policy"
  description = "Set attributes on matching routes"

  set_route_tag {
    value = 500
  }

  set_weight {
    value = 50
  }

  set_preference {
    value = 200
  }

  set_metric = 100

  set_next_hop    = "192.0.2.1"
  set_metric_type = "ospf-type2"

  set_next_hop_propagate = false
  set_multi_path         = false

  set_community {
    community = "no-export"
    criteria  = "replace"
  }

  set_additional_community = ["no-advertise", "extended:color:35"]

  set_dampening {
    half_life         = 15
    reuse_limit       = 750
    suppress_limit    = 2000
    max_suppress_time = 60
  }

  set_as_path {
    prepend {
      path_asn {
        asn   = 10
        order = 1
      }

      path_asn {
        asn   = 20
        order = 2
      }
    }

    prepend_last_as {
      asn_count = 10
    }
  }
}

data "mso_tenant_policies_set_rule_policy" "example" {
  template_id = mso_tenant_policies_set_rule_policy.example.template_id
  name        = mso_tenant_policies_set_rule_policy.example.name
}
