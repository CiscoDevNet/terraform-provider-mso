provider "mso" {
   username = "admin"
   password = "ins3965!ins3965!"
   url      = "https://173.36.219.193/"
   insecure = true
}

resource "mso_schema" "foo_schema" {
 name          = "stestnk"
 template_name = "stemplate1"
 tenant_id     = "5ea000bd2c000058f90a26ab"
}


resource "mso_schema_template_vrf" "vrf1" {
       schema_id= "${mso_schema.foo_schema.id}"
       template="stemplate1"
       name= "svrf"
       display_name="svrf"
       layer3_multicast=false
}
data "mso_schema_template_vrf" "vrf1" {
  schema_id="${mso_schema.foo_schema.id}"
  template="stemplate1"
  name= "svrf"
}
resource "mso_schema_template_vrf" "vrf" {
 count = 2
       schema_id= "${mso_schema.foo_schema.id}"
       template="${data.mso_schema_template_vrf.vrf1.template}"
       name= "svrf${count.index}"
       display_name="svrf${count.index}"
       layer3_multicast=false
}


resource "mso_schema_template_bd" "bridge_domain" {
   schema_id              = "${mso_schema.foo_schema.id}"
   template_name          = "stemplate1"
   name                   = "stestBD"
   display_name           = "stest"
   vrf_name               = "svrf"
   layer2_unknown_unicast = "proxy"
   layer2_stretch = true
}

data "mso_schema_template_bd" "bridge_domain" {
   schema_id              = "${mso_schema.foo_schema.id}"
   template_name          = "stemplate1"
   name                   = "stestBD"
}
resource "mso_schema_template_bd" "bridge_domain1" {
 count = 2
   schema_id              = "${mso_schema.foo_schema.id}"
   template_name          = "${data.mso_schema_template_bd.bridge_domain.template_name}"
   name                   = "stestBD${count.index}"
   display_name           = "stest${count.index}"
   vrf_name               = "svrf"
   layer2_unknown_unicast = "proxy"
   layer2_stretch = true
}


resource "mso_schema_template_bd_subnet" "bdsub1" {
 schema_id = "${mso_schema.foo_schema.id}"
 template_name = "stemplate1"
 bd_name = "stestBD"
 ip = "10.23.13.0/8"
 scope = "private"
 description = "Description for the subnet"
 shared = true
 no_default_gateway = false
 querier = true
}