package provider

import "testing"

func TestProviderRegistrySortsFactoriesByName(t *testing.T) {
	type testRegistrationValue struct {
		name string
	}

	registrations := []registration[testRegistrationValue]{}
	addRegistration(&registrations, "mso_z", func() testRegistrationValue { return testRegistrationValue{name: "mso_z"} })
	addRegistration(&registrations, "mso_a", func() testRegistrationValue { return testRegistrationValue{name: "mso_a"} })

	factories := registeredFactories(registrations)
	if len(factories) != 2 {
		t.Fatalf("unexpected factory count: got %d", len(factories))
	}
	if got := factories[0]().name; got != "mso_a" {
		t.Fatalf("first factory has name %q, want %q", got, "mso_a")
	}
	if got := factories[1]().name; got != "mso_z" {
		t.Fatalf("second factory has name %q, want %q", got, "mso_z")
	}
}
