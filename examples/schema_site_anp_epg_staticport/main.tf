terraform {
  required_providers {
    mso = {
      source = "CiscoDevNet/mso"
    }
  }
}

provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}
// Single Static Port
resource "mso_schema_site_anp_epg_static_port" "foo_schema_site_anp_epg_static_port" {
  schema_id            = "5c4d5bb72700000401f80948"
  site_id              = "5c7c95b25100008f01c1ee3c"
  template_name        = "Template1"
  anp_name             = "ANP"
  epg_name             = "DB"
  path_type            = "port"
  deployment_immediacy = "lazy"
  pod                  = "pod-4"
  leaf                 = "106"
  path                 = "eth1/11"
  vlan                 = 200
  micro_seg_vlan       = 3
  mode                 = "untagged"
  fex                  = "10"
}

// Create multiple Static Ports whith a list of path.
resource "mso_site" "site_test" {
  name             = var.site_name
  username         = var.site_username
  password         = var.site_password
  apic_site_id     = 105
  urls             = ["https://10.23.248.102"]
  # login_domain     = "radius_test"
  # maintenance_mode = true
  location = {
    lat  = 78.946
    long = 95.623
  }
}

resource "mso_tenant" "tenant_test" {
  name = var.tenant_name
  display_name = var.tenant_name
  site_associations {
    site_id     = mso_site.site_test.id
  }
  // site_associations = mso_site.site_test.id
}

resource "mso_schema" "schema_test" {
  name          = var.schema_name
  template_name = "Template1"
  tenant_id     = mso_tenant.tenant_test.id
}

// resource "mso_schema_template" "template" {
//   schema_id     = mso_schema.schema_test.id
//   name          = mso_schema.schema_test.template_name
//   display_name  = "Template 1"
//   tenant_id     = mso_tenant.tenant_test.id
// }

resource "mso_schema_template_vrf" "vrf" {
  schema_id     = mso_schema.schema_test.id
  template      = mso_schema.schema_test.template_name
  name          = var.vrf_name
  display_name = var.vrf_name
}

resource "mso_schema_template_bd" "bd" {
  schema_id              = mso_schema.schema_test.id
  template_name          = mso_schema.schema_test.template_name
  name                   = var.bd_name
  display_name = var.bd_name
  vrf_name               = mso_schema_template_vrf.vrf.name
  // layer2_unknown_unicast = "proxy"
  // layer2_stretch = true
}

resource "mso_schema_template_anp" "anp" {
  schema_id     = mso_schema.schema_test.id
  template      = mso_schema.schema_test.template_name
  name          = var.anp_name
  display_name = var.anp_name
}

resource "mso_schema_template_anp_epg" "db" {
  schema_id         = mso_schema.schema_test.id
  template_name     = mso_schema.schema_test.template_name
  anp_name          = mso_schema_template_anp.anp.name
  name              = var.epg_name
  display_name = var.epg_name
  bd_name           = mso_schema_template_bd.bd.name
  vrf_name          = mso_schema_template_vrf.vrf.name
}

resource "mso_schema_site" "schema_site" {
  schema_id      = mso_schema.schema_test.id
  site_id        = mso_site.site_test.id
  template_name  = mso_schema.schema_test.template_name
}

resource "mso_schema_site_anp" "anp" {
  schema_id     = mso_schema.schema_test.id
  template_name = mso_schema.schema_test.template_name
  site_id       = mso_site.site_test.id
  anp_name      = mso_schema_template_anp.anp.name
}

resource "mso_schema_site_anp_epg" "epg" {
  schema_id     = mso_schema.schema_test.id
  template_name = mso_schema.schema_test.template_name
  site_id       = mso_site.site_test.id
  anp_name      = mso_schema_site_anp.anp.anp_name
  epg_name      = mso_schema_template_anp_epg.db.name
}

resource "mso_schema_site_anp_epg_static_port" "this" {
  for_each             = toset(var.paths)
  schema_id            = mso_schema.schema_test.id
  site_id              = mso_site.site_test.id
  template_name        = mso_schema.schema_test.template_name
  anp_name             = var.anp_name
  epg_name             = mso_schema_site_anp_epg.epg.epg_name
  path_type            = var.path_type
  deployment_immediacy = var.deployment_immediacy
  pod                  = var.pod
  leaf                 = var.leaf
  path                 = each.value
  vlan                 = var.vlan
  micro_seg_vlan       = var.micro_seg_vlan
  mode                 = var.mode
  fex                  = var.fex
}