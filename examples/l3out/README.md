# L3Out with annotations and interface groups

This root configuration creates the tenant, schema and deployed VRF, fabric L3 domain, tenant policies, and L3Out template needed by `mso_l3out`. The L3Out manages both its annotation map and its interface group map. The group references a custom QoS policy, an interface routing policy, and four NetFlow monitors (`ipv4`, `ipv6`, `ce`, and `unspecified`).

Set `MSO_URL`, `MSO_USERNAME`, and `MSO_PASSWORD` for a Nexus Dashboard instance with Orchestration enabled, then supply an existing site name and the three authentication keys as Terraform input variables. For example, set `TF_VAR_site_name`, `TF_VAR_bfd_key`, `TF_VAR_bfd_multi_hop_key`, and `TF_VAR_ospf_key` in your environment before running `terraform init` and `terraform apply` in this directory. Do not put the keys in a checked-in variables file. Terraform state can contain configured keys.

Run this directory as one Terraform root configuration. The [standalone annotation](../l3out_annotation) and [standalone interface group policy](../l3out_interface_group_policy) directories demonstrate the alternative ownership model.
