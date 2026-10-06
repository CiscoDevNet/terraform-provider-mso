# Standalone L3Out annotations

This root configuration creates the tenant, deployed VRF, fabric L3 domain, L3Out template, and parent L3Out. Two `mso_l3out_annotation` resources manage separate keys. The parent L3Out leaves `annotations` omitted so the child resources own those keys independently.

Set `MSO_URL`, `MSO_USERNAME`, and `MSO_PASSWORD` for a Nexus Dashboard instance with Orchestration enabled. Set `TF_VAR_site_name` to an existing site name, then run `terraform init` and `terraform apply` in this directory. The data source reads one of the managed annotations.

See the [aggregate L3Out example](../l3out) for managing the entire annotation map through `mso_l3out.annotations`.
