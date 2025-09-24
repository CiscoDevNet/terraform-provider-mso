---
layout: "mso"
page_title: "MSO: mso_service_device_cluster"
sidebar_current: "docs-mso-data-source-service_device_cluster"
description: |-
  Data source for Service Device Cluster.
---

# mso_service_device_cluster #

Data source for a Service Device Cluster on Cisco Nexus Dashboard Orchestrator (NDO).


## GUI Information ##

* `Location` - Manage -> Service Device Template -> Devices

## Example Usage ##

```hcl
data "mso_service_device_cluster" "cluster" {
  template_id = "a1b2c3d4-e5f6-7890-1234-567890abcdef"
  name        = "my-firewall-cluster"
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the service device template where the cluster exists.
* `name` - (Required) The name of the Service Device Cluster to look up.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the Service Device Cluster.
* `id` - (Read-Only) The unique terraform identifier of the Service Device Cluster in the template.
* `description` - (Read-Only) The description of the Service Device Cluster.
* `device_mode` - (Read-Only) The operational mode of the device (e.g., layer3).
* `device_type` - (Read-Only) The type of device (e.g., firewall).
* `interface_properties` - (Read-Only) A set of interface properties associated with the cluster. Each element has the following attributes:
* `name` - (Read-Only) The name of the interface.
* `bd_uuid` - (Read-Only) The NDO UUID of the associated Bridge Domain (BD).
* `external_epg_uuid` - (Read-Only) The NDO UUID of the associated External EPG.
* `ipsla_monitoring_policy_uuid` - (Read-Only) The NDO UUID of the applied IP SLA monitoring policy.
* `qos_policy_uuid` - (Read-Only) The NDO UUID of the applied Quality of Service (QoS) policy.
* `preferred_group` - (Read-Only) Whether the interface belongs to a preferred group.
* `rewrite_source_mac` - (Read-Only) Whether source MAC address rewriting is enabled.
* `anycast` - (Read-Only) Whether anycast is enabled for this interface.
* `config_static_mac` - (Read-Only) Whether static MAC configuration is enabled.
* `is_backup_redirect_ip` - (Read-Only) Whether this is a backup redirect IP.
* `load_balance_hashing` - (Read-Only) The load balancing hashing method used.
* `pod_aware_redirection` - (Read-Only) Whether pod-aware redirection is enabled.
* `resilient_hashing` - (Read-Only) Whether resilient hashing is enabled.
* `tag_based_sorting` - (Read-Only) Whether tag-based sorting is enabled.
* `min_threshold` - (Read-Only) The minimum threshold value for redirect.
* `max_threshold` - (Read-Only) The maximum threshold value for redirect.
* `threshold_down_action` - (Read-Only) The action taken when the threshold is down.
