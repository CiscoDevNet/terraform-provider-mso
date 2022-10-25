---
layout: "mso"
page_title: "MSO: mso_schema_template_filter_entry"
sidebar_current: "docs-mso-data-source-schema_template_filter_entry"
description: |-
  MSO Schema Template Filter Entry Data Source..
---

# mso_schema_template_filter_entry #

MSO Schema Template Filter Entry Data source.

## Example Usage ##

```hcl

data "mso_schema_template_filter_entry" "filter_entry" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  name          = "Any"
  entry_name    = "Entry1"
}

```

## Argument Reference ##

* `schema_id` - (Required) SchemaID of Entry.
* `template_name` - (Required) Template name of Entry.
* `name` - (Required) Filter name of the Entry.
* `entry_name` - (Required) Name of the entry.



## Attribute Reference ##
* `display_name` - (Optional) The name of the filter as displayed on the MSO UI.
* `entry_display_name` - (Optional) The name of the entry as displayed on the MSO UI.
* `entry_description` - (Optional) The description of entry.
* `ether_type` - (Optional) The ethernet type to use for this filter entry.
* `ip_protocol` - (Optional) The IP protocol to use for this filter entry.
* `tcp_session_rules` - (Optional) A list of TCP session rules.
* `source_from` - (Optional) The source port range from.
* `source_to` - (Optional) The source port range to.
* `destination_from` - (Optional) The destination port range from.
* `destination_to` - (Optional) The destination port range to.
* `arp_flag` - (Optional) The ARP flag to use for this filter entry.
* `stateful` - (Optional) Whether this filter entry is stateful.
* `match_only_fragments` - (Optional) Whether this filter entry only matches fragments.

