package models

type Location struct {
	Lat  float64 `json:"lat,omitempty"`
	Long float64 `json:"long,omitempty"`
}

type SiteAttributes struct {
	Name           string        `json:",omitempty"`
	ApicUsername   string        `json:",omitempty"`
	ApicPassword   string        `json:",omitempty"`
	ApicSiteId     string        `json:",omitempty"`
	Labels         []interface{} `json:",omitempty"`
	Location       *Location     `json:",omitempty"`
	Url            []interface{} `json:",omitempty"`
	Platform       string        `json:",omitempty"`
	CloudProviders []interface{} `json:",omitempty"`
}

func NewSite(siteAttr SiteAttributes) *SiteAttributes {

	SiteAttributes := siteAttr
	return &SiteAttributes
}

func (siteAttributes *SiteAttributes) ToMap() (map[string]interface{}, error) {
	siteAttributeMap := make(map[string]interface{})
	A(siteAttributeMap, "name", siteAttributes.Name)
	A(siteAttributeMap, "username", siteAttributes.ApicUsername)
	A(siteAttributeMap, "password", siteAttributes.ApicPassword)
	A(siteAttributeMap, "apicSiteId", siteAttributes.ApicSiteId)
	A(siteAttributeMap, "labels", siteAttributes.Labels)
	A(siteAttributeMap, "location", siteAttributes.Location)
	A(siteAttributeMap, "urls", siteAttributes.Url)
	A(siteAttributeMap, "platform", siteAttributes.Platform)
	A(siteAttributeMap, "cloudProviders", siteAttributes.CloudProviders)

	return siteAttributeMap, nil
}