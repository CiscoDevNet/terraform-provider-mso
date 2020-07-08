provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_site" "site_test_1" {
  name             = "Cisco_MSO se"
  username         = "admin"
  password         = "noir0!234"
  apic_site_id     = "18"
  urls             = ["https://3.208.123.222"]
  # login_domain     = "radius_test"
  # maintenance_mode = true
  location = {
    lat  = 78.946
    long = 95.623
  }
}
