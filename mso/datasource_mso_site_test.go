package mso

import (
	"strings"
	"testing"
)

func TestMSOSiteDeprecationMessages(t *testing.T) {
	resourceMessage := resourceMSOSite().DeprecationMessage
	if !strings.Contains(resourceMessage, "no longer functional on Nexus Dashboard (ND) 4.0+") {
		t.Fatalf("mso_site resource deprecation message should state that it is no longer functional, got %q", resourceMessage)
	}
}
