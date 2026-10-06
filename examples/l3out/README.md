# L3Out with BGP and OSPF

This root configuration creates the tenant, schema and deployed VRF, fabric L3 domain, and L3Out template needed by `mso_l3out`. It configures both BGP and OSPF on the L3Out.

Set `MSO_URL`, `MSO_USERNAME`, and `MSO_PASSWORD` for a Nexus Dashboard instance with Orchestration enabled, then supply an existing site name through `TF_VAR_site_name` before running `terraform init` and `terraform apply` in this directory.
