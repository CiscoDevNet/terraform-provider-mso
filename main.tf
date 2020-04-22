provider "mso" {
    username = "admin"
    password = "ins3965!ins3965!"
    url = "https://173.36.219.193"
    insecure = true
}

resource "mso_schema" "foo_schema" {
 schema =  "patch-test-4",
  templates= [
    {
      name =  "Template1",
      tenantid = "5ea000bd2c000058f90a26ab",
      displayname =  "Template1",
      anps =  [],
      contracts = [],
      vrfs = [],
      bds = [],
      filters =  [],
      externalepgs =  [],
      servicegraphs = []
    }
  ],
  sites =  []
}

# resource "mso_schema_site" "schemasite1" {
#     schema = "Test1"
#     template = "Template1"
#     site = "Site1"
  
# }


