provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_label" "demolabel1" {
   label = "Demo_Label1"
   type  = "site"
  
}
