provider "mso" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.193/"
  insecure = true
}

resource "mso_schema_template_filter_entry" "filter_entry" {
    schema_id           = "5c4d5bb72700000401f80948"
	template_name       = "Template1"
	name                = "Any"
	display_name        = "Any"
	entry_name          = "testAcc"
	entry_display_name  = "testAcc"
    entry_description   = "DemoEntry"
    ether_type          = "arp"
    ip_protocol         = "eigrp"
    tcp_session_rules   = ["acknowledgement"]
	destination_from    ="unspecified"
	destination_to      ="unspecified"
	source_from         ="unspecified"
	source_to           ="unspecified"
	arp_flag            ="unspecified"
    stateful            = true
    match_only_fragments= false
}
