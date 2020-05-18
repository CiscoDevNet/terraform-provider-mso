provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
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
