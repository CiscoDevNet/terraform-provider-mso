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

resource "mso_fabric_resource_policies_port_channel_interface" "example" {
  template_id           = mso_template.fabric_resource_template.id
  name                  = "example"
  description           = "example description"
  node                  = "101"
  interfaces            = ["1/1", "1/2"]
  interface_policy_uuid = mso_fabric_policies_interface_setting.port_channel_interface.id
  interface_descriptions {
    interface   = "1/1"
    description = "1/1 description"
  }
  interface_descriptions {
    interface   = "1/2"
    description = "1/2 description"
  }
}