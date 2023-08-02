---
layout: "mso"
page_title: "MSO: mso_schema_template_filter_entry"
sidebar_current: "docs-mso-data-source-schema_template_filter_entry"
description: |-
  Data source for MSO Schema Template Filter Entry.
---

# mso_schema_template_filter_entry #

Data source for MSO Schema Template Filter Entry.

## Example Usage ##

```hcl

data "mso_schema_template_filter_entry" "example" {
  schema_id     = data.mso_schema.schema1.id
  template_name = "Template1"
  name          = "Any"
  entry_name    = "Entry1"
}

```

## Argument Reference ##

* `schema_id` - (Required) The schema ID of the Filter.
* `template_name` - (Required) The template name of the Filter.
* `name` - (Required) The name of the Filter.
* `entry_name` - (Required) The name of the Filter Entry.

## Attribute Reference ##

* `display_name` - (Read-Only) The name of the filter as displayed on the MSO UI.
* `entry_display_name` - (Read-Only) The name of the entry as displayed on the MSO UI.
* `entry_description` - (Read-Only) The description of entry.
* `ether_type` - (Read-Only) The ethernet type to use for the filter entry.
* `ip_protocol` - (Read-Only) The IP protocol to use for the filter entry.
* `tcp_session_rules` - (Read-Only) A list of TCP session rules.
* `source_from` - (Read-Only) The source port range from.
* `source_to` - (Read-Only) The source port range to.
* `destination_from` - (Read-Only) The destination port range from.
* `destination_to` - (Read-Only) The destination port range to.
* `arp_flag` - (Read-Only) The ARP flag to use for the filter entry.
* `stateful` - (Read-Only) Whether the filter entry is stateful.
* `match_only_fragments` - (Read-Only) Whether the filter entry only matches fragments.
