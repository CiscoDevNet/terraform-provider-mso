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

data "mso_site" "example" {
  name = var.site_name
}

resource "mso_tenant" "example" {
  name = "example_l3out_tenant"

  site_associations {
    site_id = data.mso_site.example.id
  }
}

resource "mso_schema" "example" {
  name = "example_l3out_schema"

  template {
    name         = "example_l3out_schema_template"
    display_name = "example_l3out_schema_template"
    tenant_id    = mso_tenant.example.id
  }
}

resource "mso_schema_site" "example" {
  schema_id     = mso_schema.example.id
  template_name = one(mso_schema.example.template).name
  site_id       = data.mso_site.example.id
}

resource "mso_schema_template_vrf" "example" {
  schema_id        = mso_schema.example.id
  template         = mso_schema_site.example.template_name
  name             = "example_l3out_vrf"
  display_name     = "example_l3out_vrf"
  layer3_multicast = true
}

resource "mso_template" "fabric_policy" {
  template_name = "example_l3out_fabric_policy"
  template_type = "fabric_policy"
  sites         = [data.mso_site.example.id]
}

resource "mso_fabric_policies_l3_domain" "example" {
  template_id = mso_template.fabric_policy.id
  name        = "example_l3_domain"
}

resource "mso_template" "tenant_policies" {
  template_name = "example_l3out_tenant_policies"
  template_type = "tenant"
  tenant_id     = mso_tenant.example.id
  sites         = [data.mso_site.example.id]
}

resource "mso_tenant_policies_match_rule_policy" "outbound" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_outbound_match"

  match_prefix {
    prefix      = "203.0.113.0/24"
    description = "Example IPv4 prefix to advertise"
  }

  match_prefix {
    prefix      = "2001:db8:10::/48"
    description = "Example IPv6 prefix to advertise"
  }
}

resource "mso_tenant_policies_route_map_policy_route_control" "outbound" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_outbound"
  description = "Permit the explicitly listed IPv4 and IPv6 prefixes for export"
  match_type  = "prefix_and_routing_policy"
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "outbound" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.outbound.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.outbound.uuid
  name                  = "advertise-example-prefixes"
  order                 = 0
  action                = "permit"
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.outbound.uuid]
}

resource "mso_tenant_policies_match_rule_policy" "interleak" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_interleak_match"

  match_prefix {
    prefix      = "10.10.0.0/16"
    description = "Example IPv4 prefix allowed for interleak"
  }

  # Prevent concurrent match-rule deletes from using stale NDO array indices.
  depends_on = [mso_tenant_policies_match_rule_policy.outbound]
}

resource "mso_tenant_policies_route_map_policy_route_control" "interleak" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_interleak"
  description = "Permit the listed prefix for interleak"
  match_type  = "routing_policy_only"

  # Prevent concurrent route-map policy deletes from using stale NDO array indices.
  depends_on = [mso_tenant_policies_route_map_policy_route_control.outbound]
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "interleak" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.interleak.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.interleak.uuid
  name                  = "interleak-example-prefix"
  order                 = 0
  action                = "permit"
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.interleak.uuid]
}

resource "mso_tenant_policies_match_rule_policy" "static_redistribution" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_static_redistribution_match"

  match_prefix {
    prefix      = "192.0.2.0/24"
    description = "Example IPv4 static route eligible for redistribution"
  }

  # Prevent concurrent match-rule deletes from using stale NDO array indices.
  depends_on = [mso_tenant_policies_match_rule_policy.interleak]
}

resource "mso_tenant_policies_route_map_policy_route_control" "static_redistribution" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_static_redistribution"
  description = "Permit the listed static route for redistribution"
  match_type  = "routing_policy_only"

  # Prevent concurrent route-map policy deletes from using stale NDO array indices.
  depends_on = [mso_tenant_policies_route_map_policy_route_control.interleak]
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "static_redistribution" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.static_redistribution.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.static_redistribution.uuid
  name                  = "redistribute-example-static-route"
  order                 = 0
  action                = "permit"
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.static_redistribution.uuid]
}

resource "mso_tenant_policies_match_rule_policy" "connected_redistribution" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_connected_redistribution_match"

  match_prefix {
    prefix      = "198.51.100.0/24"
    description = "Example IPv4 connected prefix eligible for redistribution"
  }

  # Prevent concurrent match-rule deletes from using stale NDO array indices.
  depends_on = [mso_tenant_policies_match_rule_policy.static_redistribution]
}

resource "mso_tenant_policies_route_map_policy_route_control" "connected_redistribution" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_connected_redistribution"
  description = "Permit the listed connected prefix for redistribution"
  match_type  = "routing_policy_only"

  # Prevent concurrent route-map policy deletes from using stale NDO array indices.
  depends_on = [mso_tenant_policies_route_map_policy_route_control.static_redistribution]
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "connected_redistribution" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.connected_redistribution.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.connected_redistribution.uuid
  name                  = "redistribute-example-connected-prefix"
  order                 = 0
  action                = "permit"
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.connected_redistribution.uuid]
}

resource "mso_tenant_policies_match_rule_policy" "attached_host_redistribution" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_attached_host_redistribution_match"

  match_prefix {
    prefix      = "198.51.100.10/32"
    description = "Example IPv4 attached host route eligible for redistribution"
  }

  # Prevent concurrent match-rule deletes from using stale NDO array indices.
  depends_on = [mso_tenant_policies_match_rule_policy.connected_redistribution]
}

resource "mso_tenant_policies_route_map_policy_route_control" "attached_host_redistribution" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_attached_host_redistribution"
  description = "Permit the listed attached host route for redistribution"
  match_type  = "routing_policy_only"

  # Prevent concurrent route-map policy deletes from using stale NDO array indices.
  depends_on = [mso_tenant_policies_route_map_policy_route_control.connected_redistribution]
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "attached_host_redistribution" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.attached_host_redistribution.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.attached_host_redistribution.uuid
  name                  = "redistribute-example-attached-host"
  order                 = 0
  action                = "permit"
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.attached_host_redistribution.uuid]
}

resource "mso_tenant_policies_match_rule_policy" "dampening_ipv4" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_dampening_ipv4_match"

  match_prefix {
    prefix      = "203.0.113.0/24"
    description = "IPv4-only prefix matched by the IPv4 dampening policy"
  }

  # Prevent concurrent match-rule deletes from using stale NDO array indices.
  depends_on = [mso_tenant_policies_match_rule_policy.attached_host_redistribution]
}

resource "mso_tenant_policies_set_rule_policy" "dampening_ipv4" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_dampening_ipv4_set"

  # Explicit actions avoid NDO PATCH defaults replacing omitted numeric values.
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

  set_dampening {
    half_life         = 15
    reuse_limit       = 750
    suppress_limit    = 2000
    max_suppress_time = 60
  }
}

resource "mso_tenant_policies_route_map_policy_route_control" "dampening_ipv4" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_dampening_ipv4"
  description = "Apply BGP route dampening to the matched IPv4 prefix"
  match_type  = "routing_policy_only"

  # Prevent concurrent route-map policy deletes from using stale NDO array indices.
  depends_on = [mso_tenant_policies_route_map_policy_route_control.attached_host_redistribution]
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "dampening_ipv4" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.dampening_ipv4.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.dampening_ipv4.uuid
  name                  = "dampen-example-ipv4-prefix"
  order                 = 0
  action                = "permit"
  set_rule_uuid         = mso_tenant_policies_set_rule_policy.dampening_ipv4.uuid
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.dampening_ipv4.uuid]
}

resource "mso_tenant_policies_match_rule_policy" "dampening_ipv6" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_dampening_ipv6_match"

  match_prefix {
    prefix      = "2001:db8:100::/48"
    description = "IPv6-only prefix matched by the IPv6 dampening policy"
  }

  # Prevent concurrent match-rule deletes from using stale NDO array indices.
  depends_on = [mso_tenant_policies_match_rule_policy.dampening_ipv4]
}

resource "mso_tenant_policies_set_rule_policy" "dampening_ipv6" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_dampening_ipv6_set"

  # Explicit actions avoid NDO PATCH defaults replacing omitted numeric values.
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

  set_dampening {
    half_life         = 15
    reuse_limit       = 750
    suppress_limit    = 2000
    max_suppress_time = 60
  }

  # Prevent concurrent set-rule deletes from using stale NDO array indices.
  depends_on = [mso_tenant_policies_set_rule_policy.dampening_ipv4]
}

resource "mso_tenant_policies_route_map_policy_route_control" "dampening_ipv6" {
  template_id = mso_template.tenant_policies.id
  name        = "example_l3out_dampening_ipv6"
  description = "Apply BGP route dampening to the matched IPv6 prefix"
  match_type  = "routing_policy_only"

  # Prevent concurrent route-map policy deletes from using stale NDO array indices.
  depends_on = [mso_tenant_policies_route_map_policy_route_control.dampening_ipv4]
}

resource "mso_tenant_policies_route_map_policy_route_control_context" "dampening_ipv6" {
  template_id           = mso_tenant_policies_route_map_policy_route_control.dampening_ipv6.template_id
  route_map_policy_uuid = mso_tenant_policies_route_map_policy_route_control.dampening_ipv6.uuid
  name                  = "dampen-example-ipv6-prefix"
  order                 = 0
  action                = "permit"
  set_rule_uuid         = mso_tenant_policies_set_rule_policy.dampening_ipv6.uuid
  match_rule_uuids      = [mso_tenant_policies_match_rule_policy.dampening_ipv6.uuid]
}

resource "mso_template" "l3out" {
  template_name = "example_l3out_template"
  template_type = "l3out"
  tenant_id     = mso_tenant.example.id
  sites         = [data.mso_site.example.id]
}

resource "mso_l3out" "example" {
  template_id                  = mso_template.l3out.id
  name                         = "example_l3out"
  description                  = "L3Out with BGP and OSPF"
  vrf_uuid                     = mso_schema_template_vrf.example.uuid
  l3_domain                    = mso_fabric_policies_l3_domain.example.name
  target_dscp                  = "unspecified"
  pim_enabled                  = true
  import_route_control_enabled = true
  originate_default_route      = "in_addition"

  route_maps = {
    outbound_uuid                           = mso_tenant_policies_route_map_policy_route_control_context.outbound.route_map_policy_uuid
    interleak_uuid                          = mso_tenant_policies_route_map_policy_route_control_context.interleak.route_map_policy_uuid
    static_route_redistribution_uuid        = mso_tenant_policies_route_map_policy_route_control_context.static_redistribution.route_map_policy_uuid
    connected_route_redistribution_uuid     = mso_tenant_policies_route_map_policy_route_control_context.connected_redistribution.route_map_policy_uuid
    attached_host_route_redistribution_uuid = mso_tenant_policies_route_map_policy_route_control_context.attached_host_redistribution.route_map_policy_uuid
  }

  bgp = {
    enabled                   = true
    route_dampening_ipv4_uuid = mso_tenant_policies_route_map_policy_route_control_context.dampening_ipv4.route_map_policy_uuid
    route_dampening_ipv6_uuid = mso_tenant_policies_route_map_policy_route_control_context.dampening_ipv6.route_map_policy_uuid
  }

  ospf = {
    enabled                                       = true
    area_id                                       = "0.0.0.10"
    area_type                                     = "regular"
    cost                                          = 1
    originate_default_route_always                = false
    send_redistributed_lsa                        = true
    originate_summary_lsa                         = true
    suppress_forwarding_address_in_translated_lsa = false
  }
}

data "mso_l3out" "example" {
  template_id = mso_l3out.example.template_id
  name        = mso_l3out.example.name
}
