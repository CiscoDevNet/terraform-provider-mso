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

resource "mso_tenant" "tenant1" {
  name         = "test_tenant"
  display_name = "test_tenant"
  description  = "DemoTenant"
  site_associations {
    site_id = mso_site.test_site.id
  }
}

resource "mso_schema" "schema1" {
  name          = "test_schema"
  template_name = "Template1"
  tenant_id     = mso_tenant.tenant1.id
}

resource "mso_schema_template_filter_entry" "filter_entry" {
  schema_id            = mso_schema.schema1.id
  template_name        = mso_schema.schema1.template_name
  name                 = "Filter1"
  display_name         = "Filter1"
  entry_name           = "entry1"
  entry_display_name   = "entry1"
  entry_description    = "DemoEntry"
  ether_type           = "arp"
  ip_protocol          = "eigrp"
  tcp_session_rules    = ["acknowledgement"]
  destination_from     = "unspecified"
  destination_to       = "unspecified"
  source_from          = "unspecified"
  source_to            = "unspecified"
  arp_flag             = "unspecified"
  stateful             = false
  match_only_fragments = false
}

resource "mso_schema_template_filter_entry" "filter_entry_2" {
  schema_id            = mso_schema_template_filter_entry.filter_entry.schema_id
  template_name        = mso_schema_template_filter_entry.filter_entry.template_name
  name                 = "Filter2"
  display_name         = "Filter2"
  entry_name           = "entry2"
  entry_display_name   = "entry2"
  entry_description    = "DemoEntry"
  ether_type           = "arp"
  ip_protocol          = "eigrp"
  tcp_session_rules    = ["acknowledgement"]
  destination_from     = "unspecified"
  destination_to       = "unspecified"
  source_from          = "unspecified"
  source_to            = "unspecified"
  arp_flag             = "unspecified"
  stateful             = false
  match_only_fragments = false
}

resource "mso_schema_template_contract" "template_contract" {
  schema_id     = mso_schema_template_filter_entry.filter_entry_2.schema_id
  template_name = mso_schema_template_filter_entry.filter_entry_2.template_name
  contract_name = "Contract1"
  display_name  = "Contract1"
  filter_type   = "bothWay"
  target_dscp   = "af11"
  priority      = "level1"
  scope         = "context"
  filter_relationship {
    filter_schema_id     = mso_schema_template_filter_entry.filter_entry.schema_id
    filter_template_name = mso_schema_template_filter_entry.filter_entry.template_name
    filter_name          = mso_schema_template_filter_entry.filter_entry.name
    filter_type          = "bothWay"
    action               = "permit"
    priority             = "default"
    directives           = ["log", "no_stats"]
  }
  filter_relationship {
    filter_schema_id     = mso_schema_template_filter_entry.filter_entry_2.schema_id
    filter_template_name = mso_schema_template_filter_entry.filter_entry_2.template_name
    filter_name          = mso_schema_template_filter_entry.filter_entry_2.name
    filter_type          = "bothWay"
    action               = "deny"
    priority             = "level2"
    directives           = ["log", ]
  }
}

// The below format of using filter_relationships is deprecated and might be removed in future release.
// See filter_relationship example above for new syntax.
resource "mso_schema_template_contract" "template_contract" {
  schema_id     = mso_schema_template_filter_entry.filter_entry_2.schema_id
  template_name = mso_schema_template_filter_entry.filter_entry_2.template_name
  contract_name = "Contract2"
  display_name  = "Contract2"
  filter_type   = "bothWay"
  scope         = "context"
  filter_relationships = {
    filter_schema_id     = mso_schema_template_filter_entry.filter_entry_2.schema_id
    filter_template_name = mso_schema_template_filter_entry.filter_entry_2.template_name
    filter_name          = mso_schema_template_filter_entry.filter_entry.name
  }
  directives = ["none"]
}
