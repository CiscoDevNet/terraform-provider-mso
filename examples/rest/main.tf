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

resource "mso_rest" "foo_rest" {
  path    = "api/v1/schemas/60fa24e02e000040871c2b4a"
  method  = "PATCH"
  payload = <<EOF
  [
    {
      "op": "add",
      "path": "/templates/Template4",
      "value": {
        "name"         : "Template1",
        "displayName"  : "Template1",
        "tenantId"     : "60dbb30e2e0000bfc61c2a1a"
      },
      "_updateVersion": 0
    }
  ]
  EOF
}

// resource "mso_schema_template_deploy" "template_deployer" {
//   schema_id     = "5ebb9f682c0000da45812937"
//   template_name = "Template1"
//   site_id       = "5c7c95b25100008f01c1ee3c"
//   undeploy      = true
// }
