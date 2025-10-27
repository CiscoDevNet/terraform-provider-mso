---
layout: "mso"
page_title: "MSO: mso_service_device_cluster"
sidebar_current: "docs-mso-resource-service_device_cluster"
description: |-
  Manages Service Device Clusters on Cisco Nexus Dashboard Orchestrator (NDO)
---

# mso_service_device_cluster #

Manages Service Device Clusters on Cisco Nexus Dashboard Orchestrator (NDO). This resource is supported in NDO 4.2(3) and higher.


## GUI Information ##
For ND 4.1 and later:
* `Location` - Manage -> Orchestration -> Tenant Template -> Service Device Template -> Service Device Cluster

For ND 3.2:
* `Location` - Manage -> Tenant Template -> Service Device Template -> Service Device Cluster

## Example Usage ##

```hcl
# This example creates a full stack of dependencies for a service device cluster,
# including templates, a schema, BDs, and policies.

resource "mso_template" "device_template" {
  name          = "test_device_template"
  display_name  = "test_device_template"
  template_type = "service_device"
}

resource "mso_template" "tenant_template" {
  name          = "test_tenant_template_for_device"
  display_name  = "test_tenant_template_for_device"
  template_type = "tenant_policy"
}

resource "mso_schema" "schema1" {
  name          = "SchemaForDeviceClusterTest"
  display_name  = "SchemaForDeviceClusterTest"
  template_name = "Template1"
}

resource "mso_schema_template_bd" "bd1" {
  schema_id    = mso_schema.schema1.id
  template_name = "Template1"
  name         = "test_bd_1"
  display_name = "test_bd_1"
}

resource "mso_schema_template_external_epg" "epg1" {
  schema_id    = mso_schema.schema1.id
  template_name = "Template1"
  name         = "test_epg_1"
  display_name = "test_epg_1"
}

resource "mso_tenant_policies_ipsla_monitoring_policy" "ipsla1" {
  template_id = mso_template.tenant_template.id
  name        = "test_ipsla_for_device"
  sla_type    = "icmp"
}

# Main resource for the Service Device Cluster
resource "mso_service_device_cluster" "cluster" {
  template_id = mso_template.device_template.id
  name        = "test_device_cluster"
  device_mode = "layer3"
  device_type = "firewall"

  interface_properties {
    name                         = "interface1"
    external_epg_uuid            = mso_schema_template_external_epg.epg1.uuid
    ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.ipsla1.uuid
    load_balance_hashing         = "source_ip"
    min_threshold                = 10
    max_threshold                = 90
    threshold_down_action        = "permit"
    preferred_group              = true
  }

  interface_properties {
    name                         = "interface2"
    bd_uuid                      = mso_schema_template_bd.bd1.uuid
    ipsla_monitoring_policy_uuid = mso_tenant_policies_ipsla_monitoring_policy.ipsla1.uuid
    load_balance_hashing         = "destination_ip"
  }
}
```

## Argument Reference ##

* `template_id` - (Required) The unique ID of the service device template where the cluster will be created.
* `name` - (Required) The name of the Service Device Cluster.
* `description` - (Optional) A description for the Service Device Cluster.
* `device_mode` - (Required) Specifies the operational mode of the device. Allowed values are layer1, layer2, layer3.
* `device_type` - (Required) Defines the type of device being configured. Allowed values are firewall, load_balancer, other.
* `interface_properties` - (Required) A set of interface properties blocks. The order of these blocks does not matter.
  * `name` - (Required) The name of the interface. This must be unique within the cluster.
  * `bd_uuid` - (Optional) The NDO UUID of the Bridge Domain (BD) to associate with this interface. Conflicts with external_epg_uuid.
  * `external_epg_uuid` - (Optional) The NDO UUID of the External EPG to associate with this interface. Conflicts with bd_uuid.
  * `ipsla_monitoring_policy_uuid` - (Optional) The NDO UUID of an IP SLA monitoring policy to apply to this interface.
  * `qos_policy_uuid` - (Optional) The NDO UUID of a Quality of Service (QoS) policy to apply to this interface.
  * `preferred_group` - (Optional) Whether the interface belongs to a preferred group. Defaults to false.
  * `rewrite_source_mac` - (Optional) Whether to rewrite the source MAC address. Defaults to false.
  * `anycast` - (Optional) Indicates if anycast is enabled for this interface. Defaults to false.
  * `config_static_mac` - (Optional) Indicates if static MAC configuration is enabled. Defaults to false.
  * `is_backup_redirect_ip` - (Optional) Indicates if this is a backup redirect IP. Defaults to false.
  * `load_balance_hashing` - (Optional) The load balancing hashing method. Allowed values are sourceDestinationAndProtocol, sourceIP, destinationIP. Defaults to sourceDestinationAndProtocol.
  * `pod_aware_redirection` - (Optional) Indicates if pod-aware redirection is enabled. Defaults to false.
  * `resilient_hashing` - (Optional) Indicates if resilient hashing is enabled. Defaults to false.
  * `tag_based_sorting` - (Optional) Indicates if tag-based sorting is enabled. Defaults to false.
  * `min_threshold` - (Optional) The minimum threshold value for redirect. Valid range: 0-100.
  * `max_threshold` - (Optional) The maximum threshold value for redirect. Valid range: 0-100.
  * `threshold_down_action` - (Optional) The action to take when the threshold is down. Allowed values are permit, deny, bypass. Default value is deny.

## Attribute Reference ##

* `uuid` - The NDO UUID of the Service Device Cluster.
* `id` - The unique terraform identifier of the Service Device Cluster in the template.

## Importing ##

An existing MSO Service Device Cluster can be [imported][docs-import] into this resource via its ID/path, via the following command: [docs-import]: https://www.terraform.io/docs/import/index.html


```bash
terraform import mso_service_device_cluster.cluster templateId/{template_id}/ServiceDeviceCluster/{name}
```