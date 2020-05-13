provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_rest" "sample_rest" {
    path = "/api/v1/schemas"
    payload = <<EOF
{
  "displayName": "nkauto-rest",
  "templates": [
    {
      "name": "Template1",
      "tenantId": "5e9d09482c000068500a269a",
      "displayName": "Template1",
      "anps": [],
      "contracts": [],
      "vrfs": [],
      "bds": [],
      "filters": [],
      "externalEpgs": [],
      "serviceGraphs": []
    }
  ],
  "sites": []
}
EOF
  
}
