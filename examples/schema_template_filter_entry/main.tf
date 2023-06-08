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

resource "mso_tenant" "tf_tenant" {
  name         = "tf_tenant"
  display_name = "tf_tenant"
}

resource "mso_schema" "tf_schema" {
  name = "tf_schema"
  template {
    name         = "Template1"
    display_name = "Template1"
    tenant_id    = mso_tenant.tf_tenant.id
  }
}

resource "mso_schema_template_filter_entry" "filter_entry" {
  schema_id            = mso_schema.tf_schema.id
  template_name        = "Template1"
  name                 = "Any"
  display_name         = "Any"
  entry_name           = "testAcc"
  entry_display_name   = "testAcc"
  entry_description    = "DemoEntry"
  ether_type           = "arp"
  ip_protocol          = "eigrp"
  destination_from     = "unspecified"
  destination_to       = "unspecified"
  source_from          = "unspecified"
  source_to            = "unspecified"
  arp_flag             = "unspecified"
  match_only_fragments = false
}
