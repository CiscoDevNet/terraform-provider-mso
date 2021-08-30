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

resource "mso_user" "name" {
  username      = "demoname1"
  user_password = "password12345@6"
  first_name    = "first_name"
  last_name     = "last_name"
  email         = "email@gmail.com"
  phone         = "12345678910"
  account_status= "active"
  roles{
    roleid="0000ffff0000000000000031"
    access_type = "readWrite"
  }
}
