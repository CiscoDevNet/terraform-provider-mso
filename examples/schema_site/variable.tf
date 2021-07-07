variable "site_name_1" {
  type = string
  default = "mso_site_1"
  description = "The name of the MSO site to be created"
}

variable "site_name_2" {
  type = string
  default = "mso_site_2"
  description = "The name of the MSO site to be created"
}

variable "tenant_name" {
  type = string
  default = "mso_tenant"
  description = "The name of MSO tenant to use for this configuration"
}

variable "schema_name" {
  type = string
  default = "mso_schema"
  description = "The name of the MSO schema to be created"
}

variable "template_name" {
  type = string
  default = "Template1"
  description = "Template name to which sites will be associated."
}
