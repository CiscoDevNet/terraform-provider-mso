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

resource "mso_template" "fabric_policy_template" {
  template_name = "fabric_policy_template"
  template_type = "fabric_policy"
}

resource "mso_fabric_policies_interface_setting" "port_channel_interface" {
  template_id = mso_template.fabric_policy_template.id
  type        = "portchannel"
  name        = "port_channel_interface"
}

resource "mso_template" "fabric_resource_template" {
  template_name = "fabric_resource_template"
  template_type = "fabric_resource"
}

resource "mso_fabric_resource_policies_virtual_port_channel_interface" "vpc_if" {
  template_id                 = mso_template.fabric_resource_template.id
  name                        = "tf_vpc_if"
  description                 = "Example VPC Interface"
  node_1                      = "101"
  node_2                      = "102"
  node_1_interfaces           = ["1/1", "1/10-11"]
  node_2_interfaces           = ["1/2"]
  interface_policy_group_uuid = mso_fabric_policies_interface_setting.port_channel_interface.uuid
  interface_descriptions {
    node        = "101"
    interface   = "1/1"
    description = "Interface Description 1/1"
  }
  interface_descriptions {
    node        = "102"
    interface   = "1/10"
    description = "Interface Description 1/10"
  }
}
