---
layout: "mso"
page_title: "MSO: mso_schema_site_contract_service_graph_listener"
sidebar_current: "docs-mso-data-source-schema_site_contract_service_graph_listener"
description: |-
  Data source for MSO Site Template Contract Service Graph Listener for the Azure Cloud Network Controller.
---

# mso_schema_site_contract_service_graph_listener #

Data source for MSO Site Template Contract Service Graph Listener for the Azure Cloud Network Controller.

## Example Usage ##

```hcl

data "mso_schema_site_contract_service_graph_listener" "example" {
  schema_id          = mso_schema.example.id
  template_name      = one(mso_schema.example.template).name
  contract_name      = mso_schema_template_contract.example.contract_name
  service_node_index = 0
  site_id            = mso_schema_site.example.site_id
  listener_name      = "example"
}

```

## Argument Reference ##
* `schema_id` - (Required) The schema ID of the Listener.
* `template_name` - (Required) The template name of the Listener.
* `contract_name` - (Required) The contract name of the Listener.
* `service_node_index` - (Required) The service node index of the Service Graph under which the Listener is deployed.
* `site_id` - (Required) The site ID under which the Listener is deployed.
* `listener_name` - (Required) The name of the Listener.

## Attribute Reference ##
* `protocol` - (Read-Only) The protocol of the Listener.
* `port` - (Read-Only) The port number of the Listener.
* `security_policy` - (Read-Only) The security policy of the Listener.
* `frontend_ip_dn` - (Read-Only) The frontend IP DN of the Listener.
* `ssl_certificates` - (Read-Only) The SSL Certificates information of the Listener.
  * `name` - (Read-Only) The key ring name of the SSL Certificate.
  * `target_dn` - (Read-Only) The key ring DN of the SSL Certificate.
  * `default` - (Read-Only) The default boolean flag of the SSL Certificate.
  * `certificate_store` - (Read-Only) The certificate store of the SSL Certificate.
* `rules` - (Read-Only) The Rules information of the Listener.
  * `name` - (Read-Only) The name of the Rule.
  * `floating_ip` - (Read-Only) The floating IP of the Rule.
  * `priority` - (Read-Only) The priority (index) of the Rule.
  * `host` - (Read-Only) The host of the Rule.
  * `path` - (Read-Only) The path of the Rule.
  * `action` - (Read-Only) The action of the Rule.
  * `condition` - (Read-Only) The condition of the Rule.
  * `action_type` - (Read-Only) The action type of the Rule.
  * `content_type` - (Read-Only) The content type of the Rule.
  * `port` - (Read-Only) The port number of the Rule.
  * `protocol` - (Read-Only) The protocol of the Rule.
  * `url_type` - (Read-Only) The url type of the Rule.
  * `custom_url` - (Read-Only) The custom url of the Rule.
  * `redirect_host_name` - (Read-Only) The redirect host name of the Rule.
  * `redirect_path` - (Read-Only) The redirect path of the Rule.
  * `redirect_query` - (Read-Only) The redirect query of the Rule.
  * `response_code` - (Read-Only) The response code of the Rule.
  * `response_body` - (Read-Only) The response body of the Rule.
  * `redirect_protocol` - (Read-Only) The redirect protocol of the Rule.
  * `redirect_port` - (Read-Only) The redirect port of the Rule.
  * `redirect_code` - (Read-Only) The redirect code of the Rule.
  * `target_ip_type` - (Read-Only) The target IP type of the Rule.
  * `health_check` - (Read-Only) The Health Checks information of the Listener Rule.
    * `port` - (Read-Only) The port number of the Health Checks.
    * `protocol` - (Read-Only) The port number of the Health Checks.
    * `path` - (Read-Only) The path of the Health Checks.
    * `interval` - (Read-Only) The interval(seconds) of the Health Checks.
    * `timeout` - (Read-Only) The timeout(seconds) of the Health Checks.
    * `unhealthy_threshold` - (Read-Only) The unhealthy threshold of the Health Checks.
    * `use_host_from_rule` - (Read-Only) The use host from rule of the Health Checks.
    * `success_code` - (Read-Only) The success code (code range) of the Health Checks.
    * `host` - (Read-Only) The host of the Health Checks.
  * `provider_epg_ref` - (Read-Only) The Provider EPG information of the Listener Rule.
    * `schema_id` - (Read-Only) The schema ID of the EPG Provider.
    * `template_name` - (Read-Only) The template name of the EPG Provider.
    * `anp_name` - (Read-Only) The application profile name of the EPG Provider.
    * `epg_name` - (Read-Only) The name of the EPG Provider.
