package provider

import (
	"cmp"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// registration associates a Terraform type name with its Framework factory.
// Registrations are kept separate from the provider implementation so new
// resources and data sources can register themselves without requiring a
// central list to be updated.
type registration[T any] struct {
	name    string
	factory func() T
}

var (
	resourceRegistrations   []registration[resource.Resource]
	dataSourceRegistrations []registration[datasource.DataSource]
)

func registerResource(name string, factory func() resource.Resource) {
	addRegistration(&resourceRegistrations, name, factory)
}

func registerDataSource(name string, factory func() datasource.DataSource) {
	addRegistration(&dataSourceRegistrations, name, factory)
}

func addRegistration[T any](registrations *[]registration[T], name string, factory func() T) {
	*registrations = append(*registrations, registration[T]{name: name, factory: factory})
}

func registeredFactories[T any](registered []registration[T]) []func() T {
	registrations := slices.Clone(registered)
	slices.SortFunc(registrations, func(a, b registration[T]) int {
		return cmp.Compare(a.name, b.name)
	})

	factories := make([]func() T, len(registrations))
	for index, registration := range registrations {
		factories[index] = registration.factory
	}
	return factories
}

func providerResources() []func() resource.Resource {
	return registeredFactories(resourceRegistrations)
}

func providerDataSources() []func() datasource.DataSource {
	return registeredFactories(dataSourceRegistrations)
}
