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

resource "mso_system_config" "system_config" {
  alias = "test alias"
  banner {
    message = "test message"
    state = "active"
    type = "warning"
  }
  change_control = {
    workflow = "enabled"
    number_of_approvers = 2
  }
}
