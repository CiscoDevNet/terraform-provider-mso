variable "tenant_name" {
  type        = string
  default     = "mso_tenant"
  description = "The name of MSO tenant to use for this configuration"
}

variable "site_name" {
  type        = string
  default     = "mso_site"
  description = "The name of the MSO site to be created"
}

variable "site_username" {
  type        = string
  default     = "username"
  description = "The username of the APIC"
}

variable "site_password" {
  type        = string
  default     = "password"
  description = "The password of the APIC "
}

variable "schema_name" {
  type        = string
  default     = "mso_schema"
  description = "The name of the MSO schema to be created"
}

variable "template_name" {
  type        = string
  default     = "Template1"
  description = "Template name under which you want to deploy Static Port."
}

variable "anp_name" {
  type        = string
  default     = "ANP"
  description = "ANP name under which you want to deploy Static Port."
}

variable "epg_name" {
  type        = string
  default     = "DB"
  description = "EPG name under which you want to deploy Static Port."
}

variable "vrf_name" {
  type        = string
  default     = "VRF"
  description = "VRF name under which you want to deploy ANP EPG."
}

variable "bd_name" {
  type        = string
  default     = "BD"
  description = "BD name under which you want to deploy ANP EPG."
}
