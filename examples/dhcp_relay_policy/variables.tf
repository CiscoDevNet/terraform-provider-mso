variable "tenant_name" {
  type        = string
  default     = "mso_tenant"
  description = "The name of MSO tenant to use for this configuration"
}

variable "schema_name" {
  type        = string
  default     = "mso_schema"
  description = "The name of MSO schema to use for this configuration"
}

variable "template_name" {
  type        = string
  default     = "mso_template"
  description = "The name of MSO template to use for this configuration"
}

variable "external_epg_name" {
  type        = string
  default     = "mso_external_epg"
  description = "The name of MSO external epg to use for this configuration"
}

variable "vrf_name" {
  type        = string
  default     = "mso_vrf"
  description = "The name of MSO vrf to use for this configuration"
}

variable "dhcp_server_address" {
  type        = string
  default     = "1.2.3.4"
  description = "DHCP server address."
}

variable "epg" {
  type        = string
  default     = "/schemas/621392f81d0000282a4f9d1c/templates/ACC_CREST/anps/UntitledAP1/epgs/test_epg"
  description = "EPG."
}

variable "relay_policy_name" {
  type        = string
  default     = "relay_policy"
  description = "Relay policy name."
}
