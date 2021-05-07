variable "tenant_name" {
  type = string
  default = "mso_tenant"
  description = "The name of MSO tenant to use for this configuration"
}

variable "site_name" {
  type = string
  default = "mso_site"
  description = "The name of the MSO site to be created"
}

variable "site_username" {
  type = string
  default = "username"
  description = "The username of the APIC"
}

variable "site_password" {
  type = string
  default = "password"
  description = "The password of the APIC "
}

variable "schema_name" {
  type = string
  default = "mso_schema"
  description = "The name of the MSO schema to be created"
}

variable "template_name" {
  type = string
  default = "Template1"
  description = "Template name under which you want to deploy Static Port."
}

variable "anp_name" {
  type = string
  default = "ANP"
  description = "ANP name under which you want to deploy Static Port."
}

variable "epg_name" {
  type = string
  default = "DB"
  description ="EPG name under which you want to deploy Static Port."
}

variable "vrf_name" {
  type = string
  default = "VRF"
  description = "VRF name under which you want to deploy ANP EPG."
}

variable "bd_name" {
  type = string
  default = "BD"
  description = "BD name under which you want to deploy ANP EPG."
}

variable "path_type" {
  type = string
  default = "port"
  description = "The type of the static port. Allowed values are port, vpc and dpc."
}

variable "deployment_immediacy" {
  type = string
  default = "lazy"
  description = "The deployment immediacy of the static port. Allowed values are immediate and lazy."
}

variable "pod" {
  type = string
  default = "pod-4"
  description = "The pod of the static port."
}

variable "leaf" {
  type = string
  default = "106"
  description = "The leaf of the static port."
}

variable "vlan" {
  type = number
  default = 200
  description = "The port encap VLAN id of the static port."
}

variable "micro_seg_vlan" {
  type = number
  default = 3
  description = "The microsegmentation VLAN id of the static port."
}

variable "mode" {
  type = string
  default = "untagged"
  description = "The mode of the static port. Allowed values are native, regular and untagged."
}

variable "fex" {
  type = string
  default = "10"
  description = "Fex-id to be used. This parameter will work only with the path_type as port."
}

variable "paths" {
  type    = list(string)
  description = "The path of the static port."
}