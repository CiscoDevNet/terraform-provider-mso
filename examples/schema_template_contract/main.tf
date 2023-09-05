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

data "mso_user" "user" {
  username = "admin"
}

resource "mso_tenant" "demo_tenant" {
  name         = "demo_tenant"
  display_name = "demo_tenant"
  user_associations {
    user_id = data.mso_user.user.id
  }
}

resource "mso_schema" "demo_schema" {
  name = "demo_schema"
  template {
    name         = "Template1"
    display_name = "Template1"
    tenant_id    = mso_tenant.demo_tenant.id
  }
}

resource "mso_schema_template_filter_entry" "filter_entry" {
  schema_id            = mso_schema.demo_schema.id
  template_name        = one(mso_schema.demo_schema.template).name
  name                 = "Filter1"
  display_name         = "Filter1"
  entry_name           = "entry1"
  entry_display_name   = "entry1"
}

resource "mso_schema_template_filter_entry" "filter_entry_2" {
  schema_id            = mso_schema_template_filter_entry.filter_entry.schema_id
  template_name        = mso_schema_template_filter_entry.filter_entry.template_name
  name                 = "Filter2"
  display_name         = "Filter2"
  entry_name           = "entry2"
  entry_display_name   = "entry2"
}

resource "mso_schema_template_contract" "template_contract" {
  schema_id     = mso_schema_template_filter_entry.filter_entry.schema_id
  template_name = mso_schema_template_filter_entry.filter_entry.template_name
  contract_name = "Contract1"
  display_name  = "Contract1"
  filter_type   = "bothWay"
  target_dscp   = "af11"
  priority      = "level1"
  scope         = "context"
  filter_relationship {
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
    directives           = ["log"]
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
    filter_name          = mso_schema_template_filter_entry.filter_entry_2.name
  }
  directives = ["none"]
}
