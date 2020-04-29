provider "mso" {
   username = "admin"
   password = "Crest@123456"
   # cisco-mso url
   url      = "https://10.0.9.21/"  
   insecure = true
}

resource "mso_site" "site_test_1" {
 name = "Cisco_MSO"
 username = "admin"
 password = "cisco123"
 apic_site_id = "998"
 urls = [ "https://192.168.10.102/" ]
}