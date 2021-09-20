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
  platform = "nd"
}

resource "mso_site" "test_site" {
  name         = "test_site"
  username     = "" # <APIC username>
  password     = "" # <APIC pwd>
  apic_site_id = "105"
  urls         = "" # <APIC site url>
  location = {
    lat  = 78.946
    long = 95.623
  }
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

resource "mso_schema_template_vrf" "vrf1" {
  schema_id    = mso_schema.schema1.id
  template     = mso_schema.schema1.template_name
  name         = "vrf1"
  display_name = "vrf1"
}

resource "mso_schema_template_l3out" "template_l3out" {
  schema_id         = mso_schema.schema1.id
  template_name     = mso_schema.schema1.template_name
  l3out_name        = "l3out1"
  display_name      = "l3out1"
  vrf_name          = mso_schema_template_vrf.vrf1.id
  vrf_schema_id     = mso_schema_template_vrf.vrf1.schema_id
  vrf_template_name = mso_schema_template_vrf.vrf1.template
}

resource "mso_schema_template_external_epg" "template_externalepg" {
  schema_id         = mso_schema.schema1.id
  template_name     = mso_schema.schema1.template_name
  external_epg_name = "temp_epg"
  display_name      = "temp_epg"
  vrf_name          = mso_schema_template_vrf.vrf1.id
  vrf_schema_id     = mso_schema_template_vrf.vrf1.schema_id
  vrf_template_name = mso_schema_template_vrf.vrf1.template
}

resource "mso_schema_site" "schema_site_1" {
  schema_id     = mso_schema.schema1.id
  site_id       = mso_site.test_site.id
  template_name = mso_schema.schema1.template_name
}

resource "mso_schema_site_vrf" "site_vrf" {
  template_name = mso_schema_site.schema_site_1.template_name
  site_id       = mso_schema_site.schema_site_1.site_id
  schema_id     = mso_schema_site.schema_site_1.schema_id
  vrf_name      = mso_schema_template_vrf.vrf1.name
}

resource "mso_rest" "site_l3out" {
  path    = "/mso/api/v1/schemas/${mso_schema_site.schema_site_1.schema_id}?validate=false"
  method  = "PATCH"
  payload = <<EOF
  [
    {
      "op": "add",
      "path": "/sites/${mso_schema_site.schema_site_1.site_id}-${mso_schema_site.schema_site_1.template_name}/intersiteL3outs/-",
      "value": {
        "l3outRef": {
          "l3outName": "${mso_schema_template_l3out.template_l3out.l3out_name}",
          "schemaId": "${mso_schema_site.schema_site_1.schema_id}",
          "templateName": "${mso_schema_site.schema_site_1.template_name}"
        },
      "vrfRef": {
          "schemaId": "${mso_schema_site_vrf.site_vrf.schema_id}",
          "templateName": "${mso_schema_site_vrf.site_vrf.template_name}",
          "vrfName": "${mso_schema_site_vrf.site_vrf.vrf_name}"
        }
      }
    }
  ]
  EOF
}

resource "mso_schema_site_external_epg" "site_externalepg" {
  depends_on        = [mso_rest.site_l3out]
  site_id           = mso_schema_site.schema_site_1.site_id
  schema_id         = mso_schema_site.schema_site_1.schema_id
  template_name     = mso_schema_template_external_epg.template_externalepg.template_name
  external_epg_name = mso_schema_template_external_epg.template_externalepg.external_epg_name
  l3out_name        = mso_schema_template_l3out.template_l3out.l3out_name
}
