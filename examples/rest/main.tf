provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_rest" "sample_rest" {
    path = "api/v1/schemas/5ebb9f682c0000da45812937"
    method = "PATCH"
    payload = <<EOF
[
  {
    "op": "remove",
    "path": "/templates/Template3",
    "value": {
      "name": "Template3",
      "displayName": "Templat2",
      "tenantId": "5e9d09482c000068500a269a"
    }
  }
]
EOF
  
}
