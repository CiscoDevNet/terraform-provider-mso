---
layout: "mso"
page_title: "MSO: mso_schema_template_filter_entry"
sidebar_current: "docs-mso-resource-schema_template_filter_entry"
description: |-
  Manages MSO Resource Schema Template Filter Entry
---

# mso_schema_template_filter_entry #

Manages MSO Resource Schema Template Filter Entry

## Example Usage ##

```hcl
resource "mso_schema_template_filter_entry" "filter_entry" {
		schema_id = "5c4d5bb72700000401f80948"
		template_name = "Template1"
		name = "Any"
		display_name="Any"
		entry_name = "testAcc"
		entry_display_name="testAcc"
		destination_from="unspecified"
		destination_to="unspecified"
		source_from="unspecified"
		source_to="unspecified"
		arp_flag="unspecified"
}
```

## Argument Reference ##


* `schema_id` - (Required) The schema-id where Filter entry is associated.
* `template_name` - (Required) The template associated with the filter entry.
* `name` - (Required) Filter associated with the filter entry.
* `entry_name` - (Required) The name of the entry.
* `display_name` - (Required) The name of the filter as displayed on the MSO UI.
* `entry_display_name` - (Required) The name of the entry as displayed on the MSO UI.
* `entry_description` - (Optional) Description of the entry created.
* `ether_type` - (Optional) The ethernet type to use for this filter entry. Allowed Values:  arp, fcoe, ip, ipv4, ipv6, mac-security, mpls-unicast, trill, unspecified 
* `ip_protocol` - (Optional) The IP protocol to use for this filter entry. Allowed Values:  eigrp, egp, icmp, icmpv6, igmp, igp, l2tp, ospfigp, pim, tcp, udp, unspecified 
* `tcp_session_rules` - (Optional) A list of TCP session rules. Allowed Values : acknowledgement, established, finish, synchronize, reset, unspecified 
* `source_from` - (Optional) The source port range from.
* `source_to` - (Optional) The source port range to.
* `destination_from` - (Optional) The destination port range from.
* `destination_to` - (Optional) The destination port range to.
* `arp_flag` - (Optional) The ARP flag to use for this filter entry. Allowed Values: reply, request, unspecified
* `stateful` - (Optional) Whether this filter entry is stateful. Allowed Values: true or false.
* `match_only_fragments` - (Optional) Whether this filter entry only matches fragments. Allowed Values: true or false.


## Attribute Reference ##

No attributes are exported.



