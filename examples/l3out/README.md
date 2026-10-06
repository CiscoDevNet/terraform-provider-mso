# L3Out with BGP, OSPF, and annotations

This root configuration creates the tenant, schema and deployed VRF, fabric L3 domain, and L3Out template needed by `mso_l3out`. It configures both routing protocols and manages the L3Out annotation map.

Set `MSO_URL`, `MSO_USERNAME`, and `MSO_PASSWORD` for a Nexus Dashboard instance with Orchestration enabled, then supply an existing site name through `TF_VAR_site_name` before running `terraform init` and `terraform apply` in this directory.

The [standalone annotation example](../l3out_annotation) shows how to manage individual annotation keys while leaving the parent map omitted.
