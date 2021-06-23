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
