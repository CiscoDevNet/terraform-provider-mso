terraform {
  required_providers {
    mso = {
      source = "CiscoDevNet/mso"
    }
  }
}

# Set MSO_URL, MSO_USERNAME, and MSO_PASSWORD in the environment.
provider "mso" {
  platform = "nd"
}

data "mso_site" "example" {
  name = var.site_name
}

resource "mso_tenant" "example" {
  name = "example_annotation_l3out_tenant"

  site_associations {
    site_id = data.mso_site.example.id
  }
}

resource "mso_schema" "example" {
  name = "example_annotation_l3out_schema"

  template {
    name         = "example_annotation_l3out_schema_template"
    display_name = "example_annotation_l3out_schema_template"
    tenant_id    = mso_tenant.example.id
  }
}

resource "mso_schema_site" "example" {
  schema_id     = mso_schema.example.id
  template_name = "example_annotation_l3out_schema_template"
  site_id       = data.mso_site.example.id
}

resource "mso_schema_template_vrf" "example" {
  schema_id        = mso_schema.example.id
  template         = "example_annotation_l3out_schema_template"
  name             = "example_annotation_l3out_vrf"
  display_name     = "example_annotation_l3out_vrf"
  layer3_multicast = true

  depends_on = [mso_schema_site.example]
}

resource "mso_schema_template_deploy_ndo" "vrf" {
  schema_id           = mso_schema.example.id
  template_name       = "example_annotation_l3out_schema_template"
  force_apply         = ""
  undeploy_on_destroy = true

  depends_on = [mso_schema_template_vrf.example]
}

resource "mso_template" "fabric_policy" {
  template_name = "example_annotation_l3out_fabric_policy"
  template_type = "fabric_policy"
  sites         = [data.mso_site.example.id]
}

resource "mso_fabric_policies_l3_domain" "example" {
  template_id = mso_template.fabric_policy.id
  name        = "example_annotation_l3_domain"
}

resource "mso_template" "l3out" {
  template_name = "example_annotation_l3out_template"
  template_type = "l3out"
  tenant_id     = mso_tenant.example.id
  sites         = [data.mso_site.example.id]

  depends_on = [mso_schema_template_deploy_ndo.vrf]
}

resource "mso_l3out" "example" {
  template_id = mso_template.l3out.id
  name        = "example_annotation_l3out"
  vrf_uuid    = mso_schema_template_vrf.example.uuid
  l3_domain   = mso_fabric_policies_l3_domain.example.name

  # Leave annotations omitted: each child resource owns one key.
}

resource "mso_l3out_annotation" "owner" {
  template_id = mso_template.l3out.id
  l3out_uuid  = mso_l3out.example.uuid
  key         = "owner"
  value       = "network"
}

resource "mso_l3out_annotation" "purpose" {
  template_id = mso_template.l3out.id
  l3out_uuid  = mso_l3out.example.uuid
  key         = "purpose"
  value       = "external-routing"

  depends_on = [mso_l3out_annotation.owner]
}

data "mso_l3out_annotation" "owner" {
  template_id = mso_template.l3out.id
  l3out_uuid  = mso_l3out.example.uuid
  key         = mso_l3out_annotation.owner.key
}
