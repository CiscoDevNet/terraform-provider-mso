---
layout: "mso"
page_title: "MSO: mso_schema_template_anp_epg"
sidebar_current: "docs-mso-resource-schema_template_anp_epg"
description: |-
  Manages MSO Schema Template Application Network Profiles Endpoint Groups.
---

# mso_schema_template_anp_epg #

Manages MSO Schema Template Application Network Profiles Endpoint Groups.

## Example Usage ##

```hcl

resource "mso_schema_template_anp_epg" "anp_epg" {
  schema_id                  = mso_schema.schema1.id
  template_name              = "Template1"
  anp_name                   = mso_schema_template_anp.anp1.name
  name                       = "mso_epg1"
  bd_name                    = mso_schema_template_bd.bridge_domain.name
  vrf_name                   = mso_schema_template_vrf.vrf1.name
  display_name               = "mso_epg1"
  useg_epg                   = true
  intra_epg                  = "unenforced"
  intersite_multicast_source = false
  preferred_group            = false
  epg_type                   = "service"
  access_type                = "private"
  deployment_type            = "cloud_native"
  service_type               = "azure_sql"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID under which you want to deploy Anp Epg.
* `template_name` - (Required) Template where Anp Epg to be created.
* `anp_name` - (Required) Name of the Application Network Profiles.
* `name` - (Required) Name of the Endpoint Group to manage.
* `bd_name` - (Optional) Name of the Bridge Domain. It is required when using on-premise sites.
* `bd_schema_id` - (Opional) The schemaID that defines the referenced BD.
* `bd_template_name` - (Optional) The template that defines the referenced BD.
* `vrf_name` - (Optional) Name of the Vrf. It is required when using cloud sites.
* `vrf_schema_id` - (Optional) The schemaID that defines the referenced VRF.
* `vrf_template_name` - (Optional) The template that defines the referenced VRF.
* `display_name` - (Optional) The name as displayed on the MSO web interface.
* `description` - (Optional) Description of the Anp Epg.
* `useg_epg` - (Optional) Boolean flag to enable or disable whether this is a USEG EPG. Default value is set to false.
* `intra_epg` - (Optional) Whether intra EPG isolation is enforced. choices: [ enforced, unenforced ]
* `intersite_multicast_source` - (Optional) Whether intersite multicast source is enabled. Default to false.
* `proxy_arp` - (Optional) Whether to enable Proxy ARP or not. (For Forwarding control) Default to false.
* `preferred_group` - (Optional) Boolean flag to enable or disable whether this EPG is added to preferred group.      Default value is set to false.
* `epg_type` - (Optional) EPG Type. Allowed values are `application` and `service`. Default is `application`.
* `access_type` - (Optional) Access Type of the EPG. Allowed values are `private`, `public` and `public_and_private`.
* `deployment_type` - (Optional) Deployment Type of the EPG. Allowed values are `cloud_native`, `cloud_native_managed` and `third_party`.
* `service_type` - (Optional) Service Type of the EPG. Allowed values are `azure_api_management_services`, `azure_cosmos_db`, `azure_databricks`, `azure_sql`, `azure_storage`, `azure_storage_blob`, `azure_storage_file`, `azure_storage_queue`, `azure_storage_table`, `azure_kubernetes_services`, `azure_ad_domain_services`, `azure_contain_registry`, `azure_key_vault`, `redis_cache`, `custom`.
* `custom_service_type` - (Optional) Custom Service Type of the EPG. This argument is required when `service_type` is set to `custom`.

## Attribute Reference ##

* `uuid` - (Read-Only) The NDO UUID of the EPG.

## Importing ##

An existing MSO Schema Template Application Network Profiles Endpoint Group can be [imported][docs-import] into this resource via its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_template_anp_epg.anp_epg {schema_id}/template/{template_name}/anp/{anp_name}/epg/{epg_name}
```