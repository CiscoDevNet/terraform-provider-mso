variable "site_name" {
  description = "Name of an existing site in Nexus Dashboard Orchestration."
  type        = string
}

variable "bfd_key" {
  description = "BFD authentication key for the edge interface group."
  type        = string
  sensitive   = true
}

variable "bfd_multi_hop_key" {
  description = "Multi-hop BFD authentication key for the edge interface group."
  type        = string
  sensitive   = true
}

variable "ospf_key" {
  description = "OSPF authentication key for the edge interface group."
  type        = string
  sensitive   = true
}
