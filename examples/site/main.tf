provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_site" "site_test_1" {
 name            = "Cisco_MSO_Demo"
 username        = "admin"
 password        = "ins3965!"
 apic_site_id    = "22"
 urls            = [ "https://128.107.79.157/" ]
 labels          = ["5ea7d2fc010000c57d30d9f0","5dee8ce20100001d0430d9ef"]
 location        = {
    lat = 78.946
    long = 95.623
                }
 platform        = "demo_platform"
 cloud_providers = ["demo_Cloud_provider"]

}