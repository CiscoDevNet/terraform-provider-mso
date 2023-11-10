---
layout: "mso"
page_title: "MSO: mso_schema_site_contract_service_graph_listener"
sidebar_current: "docs-mso-resource-schema_site_contract_service_graph_listener"
description: |-
  Manages MSO Site Template Contract Service Graph Listener for the Azure Cloud Network Controller.
---

# mso_schema_site_contract_service_graph_listener #

Manages MSO Site Template Contract Service Graph Listener for the Azure Cloud Network Controller.

# Note: #
This resource is only compatible with NDO versions 3.7 and 4.2+. NDO versions 4.0 and 4.1 are not supported.

## Example Usage ##

```hcl

# Sample Listener configuration for the Application Load Balancer
resource "mso_schema_site_contract_service_graph_listener" "example" {
  contract_name      = mso_schema_template_contract.example.contract_name
  listener_name      = "example"
  port               = 443
  protocol           = "https"
  schema_id          = mso_schema.example.id
  security_policy    = "default"
  service_node_index = 1
  site_id            = mso_schema_site_service_graph.example.site_id
  template_name      = one(mso_schema.example.template).name
  rules {
    action_type       = "redirect"
    content_type      = "textPlain"
    name              = "default"
    port              = 80
    priority          = 1
    protocol          = "http"
    redirect_code     = "permMoved"
    redirect_port     = 80
    redirect_protocol = "http"
    response_code     = "204"
    target_ip_type    = "unspecified"
    url_type          = "original"
    health_check {
      host                = "3.3.3.3"
      interval            = 30
      path                = "/"
      port                = 443
      protocol            = "https"
      success_code        = "200"
      timeout             = 30
      unhealthy_threshold = 3
      use_host_from_rule  = "no"
    }
  }
  ssl_certificates {
    certificate_store = "default"
    name              = "ssl_certificate_key_ring"
    # Steps to create Key Ring
    # 1. Administrative -> Security -> Certificate Authorities
    # 2. Administrative -> Security -> Key Rings
    target_dn = "uni/tn-example/certstore"
  }
}

# Sample Listener configuration for the Network Load Balancer
resource "mso_schema_site_contract_service_graph_listener" "example" {
  contract_name = mso_schema_template_contract.example.contract_name
  # Steps to configure Frontend IP Name
  # 1. Application Management -> Services -> Create a "Network Load Balancer" -> Advanced Settings -> Additional Frontend IPs -> Frontend IP Names
  frontend_ip_dn     = "uni/tn-example/clb-example/vip-2.2.2.2"
  listener_name      = "example"
  port               = 80
  protocol           = "udp"
  schema_id          = mso_schema.example.id
  security_policy    = "default"
  service_node_index = 0
  site_id            = mso_schema_site_service_graph.example.site_id
  template_name      = one(mso_schema.example.template).name
  rules {
    action_type       = "forward"
    content_type      = "textPlain"
    name              = "default"
    port              = 80
    priority          = 0
    protocol          = "udp"
    redirect_code     = "permMoved"
    redirect_port     = 0
    redirect_protocol = "inherit"
    target_ip_type    = "unspecified"
    url_type          = "original"
    health_check {
      interval            = 5
      port                = 80
      protocol            = "tcp"
      success_code        = "200-399"
      timeout             = 0
      unhealthy_threshold = 2
      use_host_from_rule  = "no"
    }
    provider_epg_ref {
      anp_name      = mso_schema_template_anp.example.name
      epg_name      = mso_schema_template_anp_epg_contract.example.epg_name
      schema_id     = mso_schema.example.id
      template_name = one(mso_schema.example.template).name
    }
  }
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the Listener.
* `template_name` - (Required) The template name of the Listener.
* `contract_name` - (Required) The contract name of the Listener.
* `site_id` - (Required) The site ID under which the Listener is deployed.
* `service_node_index` - (Required) The service node index of the Service Graph under which the Listener is deployed.
* `listener_name` - (Required) The name of the Listener.
* `protocol` - (Required) The protocol of the Listener. Allowed values are `http`, `https`, `tcp`, `udp`, `tls` and `inherit`.
* `port` - (Required) The port number of the Listener.
* `security_policy` - (Optional) The security policy of the Listener. Allowed values are `default`, `eLBSecurityPolicy-2016-08`, `eLBSecurityPolicy-FS-2018-06`, `eLBSecurityPolicy-TLS-1-2-2017-01`, `eLBSecurityPolicy-TLS-1-2-Ext-2018-06`, `eLBSecurityPolicy-TLS-1-1-2017-01`, `eLBSecurityPolicy-2015-05`, `eLBSecurityPolicy-TLS-1-0-2015-04`, `AppGwSslPolicyDefault`, `AppGwSslPolicy20150501`, `AppGwSslPolicy20170401` and `AppGwSslPolicy20170401S`.
* `frontend_ip_dn` - (Optional) The frontend IP DN of the Cloud L4-L7 Network Load Balancer device. The frontend IP should be configured in the Azure Cloud Network Controller. The frontend IP configuration option will be available under the `Advanced Settings` of the Network Load Balancer device.
* `ssl_certificates` - (Optional) The SSL Certificates information of the Listener.
  * `name` - (Required) The key ring name of the SSL Certificate.
  * `target_dn` - (Required) The key ring DN of the SSL Certificate.
  * `default` - (Read-Only) The default boolean flag of the SSL Certificate.
  * `certificate_store` - (Required) The certificate store of the SSL Certificate. Allowed values are `default`, `iam` and `acm`.
* `rules` - (Optional) The Rules information of the Listener.
  * `name` - (Required) The name of the Rule.
  * `floating_ip` - (Optional) The floating IP of the Rule.
  * `priority` - (Required) The priority (index) of the Rule.
  * `host` - (Optional) The host of the Rule.
  * `path` - (Optional) The path of the Rule.
  * `action` - (Optional) The action of the Rule.
  * `condition` - (Optional) The condition of the Rule.
  * `action_type` - (Required) The action type of the Rule. Allowed values are `fixedResponse`, `forward`, `redirect` and `haPort`.
  * `content_type` - (Optional) The content type of the Rule. Allowed values are `textPlain`, `textCSS`, `textHtml`, `appJS` and `appJson`.
  * `port` - (Required) The port number of the Rule.
  * `protocol` - (Required) The protocol of the Rule. Allowed values are `http`, `https`, `tcp`, `udp`, `tls` and `inherit`.
  * `url_type` - (Optional) The url type of the Rule. Allowed values are `original` and `custom`.
  * `custom_url` - (Optional) The custom url of the Rule.
  * `redirect_host_name` - (Optional) The redirect host name of the Rule.
  * `redirect_path` - (Optional) The redirect path of the Rule.
  * `redirect_query` - (Optional) The redirect query of the Rule.
  * `response_code` - (Optional) The response code of the Rule.
  * `response_body` - (Optional) The response body of the Rule.
  * `redirect_protocol` - (Optional) The redirect protocol of the Rule. Allowed values are `http`, `https`, `tcp`, `udp`, `tls` and `inherit`.
  * `redirect_port` - (Optional) The redirect port of the Rule.
  * `redirect_code` - (Optional) The redirect code of the Rule. Allowed values are `unknown`, `permMoved`, `found`, `seeOther` and `temporary`.
  * `target_ip_type` - (Optional) The target IP type of the Rule. Allowed values are `unspecified`, `primary` and `secondary`.
  * `health_check` - (Optional) The Health Checks information of the Listener Rule.
    * `port` - (Optional) The port number of the Health Checks.
    * `protocol` - (Optional) The port number of the Health Checks. Allowed values are `http`, `https`, `tcp`, `udp` and `tls`.
    * `path` - (Optional) The path of the Health Checks.
    * `interval` - (Optional) The interval(seconds) of the Health Checks.
    * `timeout` - (Optional) The timeout(seconds) of the Health Checks.
    * `unhealthy_threshold` - (Optional) The unhealthy threshold of the Health Checks.
    * `use_host_from_rule` - (Optional) The use host from rule of the Health Checks. Allowed values are `yes` and `no`.
    * `success_code` - (Optional) The success code (code range) of the Health Checks.
    * `host` - (Optional) The host of the Health Checks.
  * `provider_epg_ref` - (Optional) The Provider EPG information of the Listener Rule.
    * `schema_id` - (Optional) The schema ID of the EPG Provider. The `schema_id` will be used if not provided.
    * `template_name` - (Optional) The template name of the EPG Provider. The `template_name` will be used if not provided.
    * `anp_name` - (Required) The application profile name of the EPG Provider.
    * `epg_name` - (Required) The name of the EPG Provider.

## Attribute Reference ##

No attributes are exported.

## Importing ##

An existing MSO Site Template Contract Service Graph Listener can be [imported][docs-import] into this resource using its Id/path, via the following command: [docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import mso_schema_site_contract_service_graph_listener.example {schema_id}/sites/{site_id}/templates/{template_name}/contracts/{contract_name}/serviceNodes/{service_node_index}/listeners/{listener_name}
```