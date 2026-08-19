package mso

import (
	"strings"
	"testing"
)

func TestMSOSiteDeprecationMessages(t *testing.T) {
	dataSourceMessage := datasourceMSOSite().DeprecationMessage
	if !strings.Contains(dataSourceMessage, "data source is deprecated: it remains functional") {
		t.Fatalf("mso_site data source deprecation message should state that it remains functional, got %q", dataSourceMessage)
	}
	if !strings.Contains(dataSourceMessage, "mso_site resource is no longer functional") {
		t.Fatalf("mso_site data source deprecation message should identify the unusable resource, got %q", dataSourceMessage)
	}

	resourceMessage := resourceMSOSite().DeprecationMessage
	if !strings.Contains(resourceMessage, "no longer functional on Nexus Dashboard (ND) 4.0+") {
		t.Fatalf("mso_site resource deprecation message should state that it is no longer functional, got %q", resourceMessage)
	}
}
