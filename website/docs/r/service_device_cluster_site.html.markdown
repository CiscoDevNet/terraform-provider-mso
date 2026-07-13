---
layout: "mso"
page_title: "MSO: mso_service_device_cluster_site"
sidebar_current: "docs-mso-resource-service_device_cluster_site"
description: |-
  Manages the site-specific configuration of a Service Device Cluster on Cisco Nexus Dashboard Orchestrator (NDO).
---

# mso_service_device_cluster_site #

Manages the site-specific configuration of a Service Device Cluster on Cisco Nexus Dashboard Orchestrator (NDO).

This resource is supported in NDO 4.2(3) and higher.


## GUI Information ##

For ND 4.1 and later:
* `Location` - Manage -> Orchestration -> Tenant Template -> Service Device Template -> Sites -> Service Device Cluster

For ND 3.2:
* `Location` - Manage -> Tenant Template -> Service Device Template -> Sites -> Service Device Cluster


## Example Usage ##

```hcl
resource "mso_service_device_cluster_site" "cluster_site" {
  template_id = mso_template.device_template.id
  site_id     = data.mso_site.site1.id
  name        = mso_service_device_cluster.cluster.name

  domain_type = "physicalDomain"
  domain_name = mso_fabric_policies_physical_domain.physical_domain.name

  interfaces {
    name = "Internal"
    vlan = 201

    fabric_to_device_connectivity {
      pod_id    = "1"
      node_id   = ["101"]
      path      = "eth1/20"
      port_type = "port"
    }

    pbr_destinations {
      ip     = "10.0.0.10"
      mac    = "00:22:BD:F8:19:FE"
      pod_id = "1"
      tag    = "internal"
    }
  }

  interfaces {
    name = "External1"
    vlan = 202

    fabric_to_device_connectivity {
      pod_id    = "1"
      node_id   = ["101"]
      path      = "eth1/25"
      port_type = "port"
    }
  }
}
```


## Argument Reference ##

* `template_id` - (Required) The ID of the service device template that contains the Service Device Cluster. Changing this forces a new resource.
* `name` - (Required) The name of the Service Device Cluster to configure on the site. Must match the `name` of the corresponding `mso_service_device_cluster` resource. Changing this forces a new resource.
* `site_id` - (Required) The ID of the site on which to configure the Service Device Cluster. Changing this forces a new resource.

* `domain_type` - (Optional) The type of domain associated with the Service Device Cluster on the site. Allowed values are `physicalDomain`, `vmmDomain`. Must be used together with `domain_name` and cannot be combined with `domain_dn`.
* `vmm_domain_type` - (Optional) The VMM domain provider type. Required when `domain_type` is `vmmDomain` and must not be set when `domain_type` is `physicalDomain`. Allowed values are `VMware`, `Microsoft`, `Redhat`.
* `domain_name` - (Optional) The name of the domain associated with the Service Device Cluster on the site. Must be used together with `domain_type` and cannot be combined with `domain_dn`.
* `domain_dn` - (Optional) The distinguished name of the domain associated with the Service Device Cluster on the site. Must start with `uni/phys-` for a physical domain or `uni/vmmp-` for a VMM domain. Cannot be combined with `domain_type`, `vmm_domain_type`, or `domain_name`.

* `high_availability_mode` - (Optional) The high availability mode of the Service Device Cluster on the site. Allowed values are `activeActive`, `activeStandby`, `notAvailable`. Defaults to `notAvailable`. Changing this forces a new resource, because transitioning between `activeActive` (per-interface domains) and the other modes (device-level domain) is a structural change to the device entry on NDO.
* `promiscuous_mode` - (Optional) Whether promiscuous mode is enabled on the Service Device Cluster on the site.
* `trunking_port` - (Optional) Whether the Service Device Cluster on the site uses a trunking port.
* `vlan` - (Optional) The device-level VLAN ID. Valid range: 1-4094. Set this only when `high_availability_mode` is `activeStandby`; for other modes use the interface-level `vlan` (regular L3 devices) or the per-path `vlan` inside `fabric_to_device_connectivity` (`activeActive`).

* `interfaces` - (Required) An ordered list of interface blocks describing the per-site configuration of every cluster interface. Must contain at least one entry, and the `name` values must match the interfaces declared in `mso_service_device_cluster.interface_properties`. Adding or removing entries forces a new resource (NDO server-side validation requires the device entry to be torn down and rebuilt when the interface set is reshaped). Content edits within existing interfaces — including the per-interface `vlan`, `fabric_to_device_connectivity`, `vm_information`, `enhanced_lag_policy`, and `pbr_destinations` — are submitted as an in-place wholesale `/interfaces` replace through the Update path. Keep entries in a stable order across applies to avoid unnecessary churn.
  * `name` - (Required) The name of the interface.
  * `vlan` - (Optional) The VLAN ID of the interface. Valid range: 1-4094. Must not be set when the matching cluster `interface_properties` binds to an `external_epg_uuid` (L3out interface); NDO rejects a VLAN on L3out interfaces.
  * `domain_name` - (Optional) The name of the physical domain associated with the interface. Mutually exclusive with `domain_dn`. Only valid when `high_availability_mode` is `activeActive`; in that mode every interface must configure its own physical domain (the device-level domain attributes are derived from the first interface and any device-level values in config are ignored). Only physical domains are supported at interface scope.
  * `domain_dn` - (Optional) The distinguished name of the physical domain associated with the interface. Must start with `uni/phys-`. Mutually exclusive with `domain_name`. Only valid when `high_availability_mode` is `activeActive`.
  * `fabric_to_device_connectivity` - (Optional) A list of fabric-to-device connectivity paths for the interface. Allowed only when the device uses a physical domain. Mutually exclusive with `vm_information`.
      * `pod_id` - (Required) The pod ID of the fabric path.
      * `node_id` - (Required) The node ID(s) of the fabric path, as a list of strings. Provide a single element for `port_type` `port` and `dpc`. Provide exactly two elements for `port_type` `vpc`.
      * `path` - (Required) The path on the node. For `port_type` `port` this is the interface (e.g. `eth1/1`). For `port_type` `dpc` and `vpc` this is the policy group name.
      * `port_type` - (Required) The type of port used for the fabric path. Allowed values are `port`, `vpc`, `dpc`.
      * `tag` - (Optional) The tag of the fabric path.
      * `vlan` - (Optional) The VLAN ID carried on this fabric path. Valid range: 1-4094. Used when `high_availability_mode` is `activeActive`, where each fabric path carries its own access VLAN.
  * `vm_information` - (Optional) A list of VM information entries for the interface. Allowed only when the device uses a VMM domain. Mutually exclusive with `fabric_to_device_connectivity`.
      * `vm_name` - (Required) The name of the VM.
      * `vnic_name` - (Required) The name of the vNIC on the VM.
      * `pod_id` - (Optional) The pod ID of the fabric path the VM interface attaches to.
      * `node_id` - (Optional) The node ID(s) of the fabric path the VM interface attaches to, as a list of strings. Provide a single element for `port_type` `port` and `dpc`. Provide exactly two elements for `port_type` `vpc`.
      * `path` - (Optional) The path on the node the VM interface attaches to. For `port_type` `port` this is the interface (e.g. `eth1/1`). For `port_type` `dpc` and `vpc` this is the policy group name.
      * `port_type` - (Optional) The type of port used for the VM interface's fabric path. Allowed values are `port`, `vpc`, `dpc`.
  * `enhanced_lag_policy` - (Optional) The name of the enhanced LAG policy associated with the interface. Only valid when the device uses a VMM domain.
  * `pbr_destinations` - (Optional) A list of policy-based redirect (PBR) destinations for the interface.
      * `ip` - (Optional) The IP address of the PBR destination. Required for L3 device clusters; omit for L1 clusters with `high_availability_mode` `activeActive` or `activeStandby`, which carry only `mac` and `tag`.
      * `mac` - (Optional) The MAC address of the PBR destination.
      * `pod_id` - (Optional) The pod ID of the PBR destination.
      * `additional_tracking_ip` - (Optional) The additional IP address used for tracking the PBR destination. NDO defaults this to `0.0.0.0` for L3 PBR destinations when omitted, and the provider treats that default as equivalent to leaving the attribute unset (no drift).
      * `weight` - (Optional) The weight of the PBR destination. Valid range: 1-10.
      * `is_backup` - (Optional) Whether the PBR destination is a backup destination.
      * `tag` - (Optional) The tag of the PBR destination.


## Attribute Reference ##

* `id` - The unique terraform identifier of the Service Device Cluster site configuration, in the form `templateId/{template_id}/site/{site_id}/ServiceDeviceCluster/{name}`.


## Importing ##

An existing MSO Service Device Cluster site configuration can be [imported][docs-import] into this resource via its ID, using the following command: [docs-import]: https://www.terraform.io/docs/import/index.html


```bash
terraform import mso_service_device_cluster_site.cluster_site templateId/{template_id}/site/{site_id}/ServiceDeviceCluster/{name}
```
