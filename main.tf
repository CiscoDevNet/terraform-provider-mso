provider "mso" {
    username = "admin"
    password = "ins3965!ins3965!"
    url = "https://173.36.219.193"
    insecure = true
}

resource "mso_schema" "schema1" {
    schema = "Test1"
  
}

resource "mso_schema_site" "schemasite1" {
    schema = "Test1"
    template = "Template1"
    site = "Site1"
  
}


