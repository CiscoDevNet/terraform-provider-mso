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
    name         = "tf_template"
    display_name = "tf_template"
    tenant_id    = mso_tenant.tf_tenant.id
  }
}

resource "mso_schema_template_vrf" "tf_vrf" {
  name         = "tf_vrf"
  display_name = "tf_vrf"
  schema_id    = mso_schema.tf_schema.id
  template     = "tf_template"
}

resource "mso_schema_template_bd" "tf_bd" {
  schema_id              = mso_schema.tf_schema.id
  template_name          = "tf_template"
  name                   = "tf_bd"
  display_name           = "tf_bd"
  layer2_unknown_unicast = "proxy"
  vrf_name               = mso_schema_template_vrf.tf_vrf.name
}

resource "mso_template" "tf_tenant_template" {
  template_name = "tf_tenant_template"
  template_type = "tenant"
  tenant_id     = mso_tenant.tf_tenant.id
}

resource "mso_tenant_policies_endpoint_mac_tag_policy" "endpoint_mac_tag_bd" {
  template_id = mso_template.tf_tenant_template.id
  mac         = "AA:BB:A1:B2:C3:D4"
  bd_uuid     = mso_schema_template_bd.tf_bd.uuid

  tag_annotations {
    key   = "annotation_key_1"
    value = "annotation_value_1"
  }

  tag_annotations {
    key   = "annotation_key_2"
    value = "annotation_value_2"
  }

  policy_tags {
    key   = "policy_key_1"
    value = "policy_value_1"
  }

  policy_tags {
    key   = "policy_key_2"
    value = "policy_value_2"
  }
}

resource "mso_tenant_policies_endpoint_mac_tag_policy" "endpoint_mac_tag_vrf" {
  template_id = mso_template.tf_tenant_template.id
  mac         = "AA:BB:A1:B2:C3:D4"
  vrf_uuid    = mso_schema_template_bd.tf_vrf.uuid

  tag_annotations {
    key   = "annotation_key_1"
    value = "annotation_value_1"
  }

  tag_annotations {
    key   = "annotation_key_2"
    value = "annotation_value_2"
  }

  policy_tags {
    key   = "policy_key_1"
    value = "policy_value_1"
  }

  policy_tags {
    key   = "policy_key_2"
    value = "policy_value_2"
  }
}
