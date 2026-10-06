# Standalone L3Out interface group policy

This root configuration creates the tenant, deployed VRF, fabric L3 domain, L3Out template, and tenant policy references needed for a standalone `mso_l3out_interface_group_policy`. The parent L3Out leaves `interface_groups` omitted. The policy uses custom QoS, interface routing, all four NetFlow monitor reference types, and BFD, multi-hop BFD, and OSPF authentication.

Set `MSO_URL`, `MSO_USERNAME`, and `MSO_PASSWORD` for a Nexus Dashboard instance with Orchestration enabled. Supply an existing site name and three authentication keys through `TF_VAR_site_name`, `TF_VAR_bfd_key`, `TF_VAR_bfd_multi_hop_key`, and `TF_VAR_ospf_key`, then run `terraform init` and `terraform apply` in this directory. Do not put the keys in a checked-in variables file. Terraform state can contain configured keys; the data source cannot read them back from Orchestration.
