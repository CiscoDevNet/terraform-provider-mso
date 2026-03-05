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

resource "mso_fabric_policies_interface_setting" "physical_interface" {
	template_id = mso_template.fabric_policy_template.id
	type        = "physical"
	name        = "physical_interface"
}

resource "mso_template" "fabric_resource_template" {
	template_name = "fabric_resource_template"
	template_type = "fabric_resource"
}

resource "mso_fabric_resource_policies_physical_interface" "physical" {
  template_id           = mso_template.fabric_resource_template.id
  name                  = "physical_interface"
  description           = "Terraform test Physical Interface"
  nodes                 = ["101", "102"]
  interfaces            = ["1/1", "1/2"]
  interface_policy_uuid = mso_fabric_policies_interface_setting.physical_interface.id
  interface_descriptions {
    interface   = "1/1"
    description = "Interface Description 1/1"
  }
  interface_descriptions {
    interface   = "1/2"
    description = "Interface Description 1/2"
  }
}

resource "mso_fabric_resource_policies_physical_interface" "breakout" {
  template_id   = mso_template.fabric_resource_template.id
  name          = "breakout_mode_physical_interface"
  description   = "Terraform test Physical Interface with breakout mode"
  nodes         = ["101", "102"]
  interfaces    = ["1/1", "1/2"]
  breakout_mode = "4x100G"
  interface_descriptions {
    interface   = "1/1"
    description = "Interface Description 1/1"
  }
  interface_descriptions {
    interface   = "1/2"
    description = "Interface Description 1/2"
  }
}