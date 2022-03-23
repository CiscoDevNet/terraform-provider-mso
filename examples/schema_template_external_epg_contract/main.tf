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

resource "mso_tenant" "test_tenant" {
  name         = "eepg_contract_tenant"
  display_name = "eepg_contract_tenant"
  description  = "DemoTenant"
}

resource "mso_schema" "test_schema" {
  name          = "eepg_contract_schema"
  template_name = "eepg_contract_template"
  tenant_id     = mso_tenant.test_tenant.id
}

resource "mso_schema_template_vrf" "test_vrf" {
  schema_id        = mso_schema.test_schema.id
  template         = mso_schema.test_schema.template_name
  name             = "eepg_contract_vrf"
  display_name     = "eepg_contract_vrf"
  layer3_multicast = false
}

resource "mso_schema_template_external_epg" "template_externalepg" {
  schema_id         = mso_schema.test_schema.id
  template_name     = mso_schema.test_schema.template_name
  external_epg_name = "eepg"
  display_name      = "eepg"
  vrf_name          = mso_schema_template_vrf.test_vrf.name
  external_epg_type = "on-premise"
}

resource "mso_schema_template_filter_entry" "filter_entry" {
  schema_id            = mso_schema.test_schema.id
  template_name        = mso_schema.test_schema.template_name
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

resource "mso_schema_template_contract" "template_contract" {
  schema_id     = mso_schema.test_schema.id
  template_name = mso_schema.test_schema.template_name
  contract_name = "Contract1"
  display_name  = "Contract1"
  filter_type   = "bothWay"
  scope         = "context"
  filter_relationship {
    filter_schema_id     = mso_schema_template_filter_entry.filter_entry.schema_id
    filter_template_name = mso_schema_template_filter_entry.filter_entry.template_name
    filter_name          = mso_schema_template_filter_entry.filter_entry.name
  }
  directives = ["none"]
}

resource "mso_schema_template_contract" "template_contract_2" {
  schema_id     = mso_schema.test_schema.id
  template_name = mso_schema.test_schema.template_name
  contract_name = "Contract2"
  display_name  = "Contract2"
  filter_type   = "bothWay"
  scope         = "context"
  filter_relationship {
    filter_schema_id     = mso_schema_template_filter_entry.filter_entry.schema_id
    filter_template_name = mso_schema_template_filter_entry.filter_entry.template_name
    filter_name          = mso_schema_template_filter_entry.filter_entry.name
  }
  directives = ["none"]
}

resource "mso_schema_template_external_epg_contract" "consumer_contract" {
  schema_id                 = mso_schema.test_schema.id
  template_name             = mso_schema.test_schema.template_name
  contract_name             = mso_schema_template_contract.template_contract.contract_name
  external_epg_name         = mso_schema_template_external_epg.template_externalepg.external_epg_name
  relationship_type         = "consumer"
  contract_schema_id        = mso_schema_template_contract.template_contract.schema_id
  contract_template_name    = mso_schema_template_contract.template_contract.template_name
}

resource "mso_schema_template_external_epg_contract" "provider_contract" {
  schema_id                 = mso_schema.test_schema.id
  template_name             = mso_schema.test_schema.template_name
  contract_name             = mso_schema_template_contract.template_contract.contract_name
  external_epg_name         = mso_schema_template_external_epg.template_externalepg.external_epg_name
  relationship_type         = "provider"
  contract_schema_id        = mso_schema_template_contract.template_contract.schema_id
  contract_template_name    = mso_schema_template_contract.template_contract.template_name
}

resource "mso_schema_template_external_epg_contract" "provider_contract_2" {
  schema_id                 = mso_schema.test_schema.id
  template_name             = mso_schema.test_schema.template_name
  contract_name             = mso_schema_template_contract.template_contract_2.contract_name
  external_epg_name         = mso_schema_template_external_epg.template_externalepg.external_epg_name
  relationship_type         = "provider"
  contract_schema_id        = mso_schema_template_contract.template_contract_2.schema_id
  contract_template_name    = mso_schema_template_contract.template_contract_2.template_name
}
