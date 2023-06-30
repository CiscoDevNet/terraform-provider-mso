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

# remote location with password authentication
resource "mso_remote_location" "password" {
  name        = "remote_location_name"
  description = "remote location with password authentication"
  protocol    = "scp"
  hostname    = "10.0.0.1"
  path        = "/tmp"
  username    = "admin"
  password    = "password"
}

# remote location with ssh key authentication
resource "mso_remote_location" "ssh" {
  name        = "remote_location_name"
  description = "remote location with ssh key authentication"
  protocol    = "scp"
  hostname    = "10.0.0.2"
  path        = "/tmp"
  username    = "admin"
  ssh_key     = "ssh_key"
  passphrase  = "passphrase"
}

