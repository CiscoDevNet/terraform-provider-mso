## 0.7.0 (August 26, 2022)
IMPROVEMENTS:
- Ability to add Microsoft and Redhat domains in mso_schema_site_anp_epg_domain (deprecate dn attribute) (#130)
- Ability to add multiple templates using the resource mso_schema (#135)
- Add import support for mso_schema_template_vrf_contract
- Apple M1 support
- Add name attribute to the subnet in resource_mso_schema_site_vrf_region (#140)

BUG FIXES:
- Fix import in mso_schema_template_anp_epg which caused configured values to be null (#129)
- Fix subnet name idempotency issue in mso_schema_template_external_epg_subnet
- Update example and update documentation for mso_schema_template_external_epg_subnet
- Fix mso_schema_template_anp_epg if only bd or vrf is provided
- Fix terraform import of schema_template_external_epg_contract
- Fix mso_schema_template_bd import issues (#136)
- Fix mso_schema_template_external_epg_subnet recreation issue when IP was changed
- Fix bug issue for mso_site import when platform is nd to keep consistency as mso platform (#137)
- Fix mso_schema_template_external_epg_contract relation type and update example

## 0.6.0 (March 15, 2022)
IMPROVEMENTS:
- Add Service EPG support for cloud sites (#113)

BUG FIXES:
- Update mso_schema_site_anp_epg_static_port documentation and example for path_type = vpc/dpc (#118)

## 0.5.0 (February 28, 2022)
IMPROVEMENTS:
- Add arp_flooding, virtual_mac_address, unicast_routing, ipv6_unknown_multicast_flooding, multi_destination_flooding and unknown_multicast_flooding in mso_schema_template_bd

BUG FIXES:
- Fix import documentation for mso_schema_template_external_epg_subnet, mso_schema_site_bd_subnet and mso_schema_template_filter_entry.
- Fix aci_site import and idempotency issues with NDO
- Add check to ignore error code 141 Resource Not Present to all Delete methods
- Fix idempotency issue of mso_schema_template_bd

## 0.4.1 (December 17, 2021)
BUG FIXES:
- Fix documentation for mso_schema_template_external_epg, mso_schema_template_external_epg_subnet and mso_schema_template_contract.

## 0.4.0 (December 16, 2021)
BEHAVIOR CHANGE:
- Modify mso_schema_template_deploy so destroy undeploy all sites

BUG FIXES:
- Fix idemptotency issue with site_id in template_external_epg

## 0.3.4 (December 16, 2021)

BUG FIXES:
- Fix login payload issue due to NDO API change (mso-go-client v1.2.6 update)
- Fix API change in l3outRef and add some error catching conditions

## 0.3.3 (December 15, 2021)

BUG FIXES:
- Fix login domain issue with NDO (mso-go-client v1.2.3 update)
- Fix mso_user data source when running on NDO
- Fix mso_schema_template_external_epg and mso_schema_site_external_epg when running on NDO

## 0.3.2 (October 29, 2021)

BUG FIXES:
- Fix issue in mso_tenant with Cloud Accounts

## 0.3.1 (September 30, 2021)

BUG FIXES:
- Fix DeletebyId crash issue with nil pointer (mso-go-client v1.2.2 update)
- Make zone attribute not mandatory in mso_schema_site_vrf_region and mso_schema_site_vrf_region_cidr_subnet (to support Azure sites)

## 0.3.0 (September 24, 2021)

BUG FIXES:
- Fix mso_schema_site_anp_epg Import when VRF or BD is not configured

IMPROVEMENTS:
- Added new resource mso_schema_site_external_epg
- Updated mso_user to support MSO/ND platform
- Updated resource and  datasource mso_site to work with ND-based MSO / NDO
- Updated mso_schema_template_contract to support multiple filter_relationship and deprecate filter_relationships attribute.

## 0.2.1 (July 20, 2021)

BUG FIXES:
- Added examples and documentations for Multi-Site associations
- Fix ressources creation and update when PATCH request return 204 No Content

## 0.2.0 (June 23, 2021)

BUG FIXES:
- Fix mso_schema_site_anp_epg_static_leaf documentation example
- Fix mso_site documentation for location attribute (locations -> location)

IMPROVEMENTS:
- Support of ND authentication (platform = "nd") for MSO v3.2+
- Updated mso_schema_site and mso_schema_template_vrf resource docs name
- Updated all terraform examples with required_provider section added in recent Terraform versions
- Added example of how to use for_each to support multiple static ports
- Added import capability on all resources
- Add support for fex ports in mso_schema_site_anp_epg_static_port
- Cleanup of go.mod and vendor directory

## 0.1.5 (January 25, 2021)

BUG FIXES:
- Fixed an issue with mso_tenant resource, where users were not able to attach the On-Prem site to the tenants due to panic error.

## 0.1.3 (October 28, 2020)

IMPROVEMENTS:

- Enabled the Auth-token sharing between API calls.

BUG FIXES:
- Fixed an issue with mso_schema_template_contract resource, where wrong filter was being read when there are multiple filters added to the same contract.

## 0.1.2 (September 3, 2020)

IMPROVEMENTS:

- Renamed resources for naming consistency.
- First Terraform Registry release.

## 0.1.1 (July 22, 2020)

IMPROVEMENTS:

- Added new resources to manage selectors on various levels.
- Added support for login_domains to site resource.
- Added Azure and AWS account support.
- Renamed resources for better formatting.

## 0.1.0 (June 17, 2020)

- Initial Release
