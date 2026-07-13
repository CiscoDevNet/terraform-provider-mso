---
layout: "mso"
page_title: "MSO: mso_service_device_cluster_site"
sidebar_current: "docs-mso-data-source-service_device_cluster_site"
description: |-
  Data source for the site-specific configuration of a Service Device Cluster.
---

# mso_service_device_cluster_site #

Data source for the site-specific configuration of a Service Device Cluster on Cisco Nexus Dashboard Orchestrator (NDO).

This data source is supported in NDO 4.2(3) and higher.


## GUI Information ##

For ND 4.1 and later:
* `Location` - Manage -> Orchestration -> Tenant Template -> Service Device Template -> Sites -> Service Device Cluster

For ND 3.2:
* `Location` - Manage -> Tenant Template -> Service Device Template -> Sites -> Service Device Cluster


## Example Usage ##

```hcl
data "mso_service_device_cluster_site" "cluster_site" {
  template_id = mso_template.device_template.id
  site_id     = data.mso_site.site1.id
  name        = mso_service_device_cluster.cluster.name
}
```


## Argument Reference ##

* `template_id` - (Required) The ID of the service device template that contains the Service Device Cluster.
* `site_id` - (Required) The ID of the site on which the Service Device Cluster is configured.
* `name` - (Required) The name of the Service Device Cluster configured on the site.


## Attribute Reference ##

* `id` - (Read-Only) The unique terraform identifier of the Service Device Cluster site configuration.
* `domain_type` - (Read-Only) The type of domain associated with the Service Device Cluster on the site.
* `vmm_domain_type` - (Read-Only) The VMM domain provider type when `domain_type` is `vmmDomain`.
* `domain_name` - (Read-Only) The name of the domain associated with the Service Device Cluster on the site.
* `domain_dn` - (Read-Only) The distinguished name of the domain associated with the Service Device Cluster on the site.
* `high_availability_mode` - (Read-Only) The high availability mode of the Service Device Cluster on the site.
* `promiscuous_mode` - (Read-Only) Whether promiscuous mode is enabled on the Service Device Cluster on the site.
* `trunking_port` - (Read-Only) Whether the Service Device Cluster on the site uses a trunking port.
* `vlan` - (Read-Only) The device-level VLAN ID of the Service Device Cluster on the site.
* `interfaces` - (Read-Only) A list of interface entries describing the per-site configuration of every cluster interface. Each element has the following attributes:
  * `name` - (Read-Only) The name of the interface.
  * `vlan` - (Read-Only) The VLAN ID of the interface.
  * `domain_name` - (Read-Only) The name of the physical domain associated with the interface (populated when `high_availability_mode` is `activeActive`; only physical domains are supported at interface scope).
  * `domain_dn` - (Read-Only) The distinguished name of the physical domain associated with the interface (populated when `high_availability_mode` is `activeActive`).
  * `fabric_to_device_connectivity` - (Read-Only) A list of fabric-to-device connectivity paths for the interface. Each element has the following attributes:
      * `pod_id` - (Read-Only) The pod ID of the fabric path.
      * `node_id` - (Read-Only) The node ID(s) of the fabric path, as a list of strings. A single element for `port_type` `port` and `dpc`, two elements for `port_type` `vpc`.
      * `path` - (Read-Only) The path on the node.
      * `port_type` - (Read-Only) The type of port used for the fabric path.
      * `tag` - (Read-Only) The tag of the fabric path.
      * `vlan` - (Read-Only) The VLAN ID carried on this fabric path (used when `high_availability_mode` is `activeActive`).
  * `vm_information` - (Read-Only) A list of VM information entries for the interface. Each element has the following attributes:
      * `vm_name` - (Read-Only) The name of the VM.
      * `vnic_name` - (Read-Only) The name of the vNIC on the VM.
      * `pod_id` - (Read-Only) The pod ID of the fabric path the VM interface attaches to.
      * `node_id` - (Read-Only) The node ID(s) of the fabric path the VM interface attaches to.
      * `path` - (Read-Only) The path on the node the VM interface attaches to.
      * `port_type` - (Read-Only) The type of port used for the VM interface's fabric path.
  * `enhanced_lag_policy` - (Read-Only) The name of the enhanced LAG policy associated with the interface.
  * `pbr_destinations` - (Read-Only) A list of policy-based redirect (PBR) destinations for the interface. Each element has the following attributes:
      * `ip` - (Read-Only) The IP address of the PBR destination.
      * `mac` - (Read-Only) The MAC address of the PBR destination.
      * `pod_id` - (Read-Only) The pod ID of the PBR destination.
      * `additional_tracking_ip` - (Read-Only) The additional IP address used for tracking the PBR destination.
      * `weight` - (Read-Only) The weight of the PBR destination.
      * `is_backup` - (Read-Only) Whether the PBR destination is a backup destination.
      * `tag` - (Read-Only) The tag of the PBR destination.
