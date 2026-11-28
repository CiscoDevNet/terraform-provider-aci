package provider

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestProviderRegistryMatchesProvider(t *testing.T) {
	provider := &AciProvider{}

	registeredResources := resourceFactoryNames(t, providerResources())
	providerResourceNames := resourceFactoryNames(t, provider.Resources(context.Background()))
	if !reflect.DeepEqual(registeredResources, providerResourceNames) {
		t.Fatalf("registered resources %v do not match provider resources %v", registeredResources, providerResourceNames)
	}

	registeredDataSources := dataSourceFactoryNames(t, providerDataSources())
	providerDataSourceNames := dataSourceFactoryNames(t, provider.DataSources(context.Background()))
	if !reflect.DeepEqual(registeredDataSources, providerDataSourceNames) {
		t.Fatalf("registered data sources %v do not match provider data sources %v", registeredDataSources, providerDataSourceNames)
	}

	registeredFunctions := functionFactoryNames(t, providerFunctions())
	providerFunctionNames := functionFactoryNames(t, provider.Functions(context.Background()))
	if !reflect.DeepEqual(registeredFunctions, providerFunctionNames) {
		t.Fatalf("registered functions %v do not match provider functions %v", registeredFunctions, providerFunctionNames)
	}
}

func TestProviderRegistryEntries(t *testing.T) {
	resourceNames := make(map[string]struct{}, len(resourceRegistrations))
	for _, registration := range resourceRegistrations {
		if registration.factory == nil {
			t.Fatalf("resource %q has a nil factory", registration.name)
		}
		if _, exists := resourceNames[registration.name]; exists {
			t.Fatalf("resource %q is registered more than once", registration.name)
		}
		resourceNames[registration.name] = struct{}{}

		response := &resource.MetadataResponse{}
		registration.factory().Metadata(
			context.Background(),
			resource.MetadataRequest{ProviderTypeName: "aci"},
			response,
		)
		if response.TypeName != registration.name {
			t.Fatalf("resource registration %q has metadata type name %q", registration.name, response.TypeName)
		}
	}

	dataSourceNames := make(map[string]struct{}, len(dataSourceRegistrations))
	for _, registration := range dataSourceRegistrations {
		if registration.factory == nil {
			t.Fatalf("data source %q has a nil factory", registration.name)
		}
		if _, exists := dataSourceNames[registration.name]; exists {
			t.Fatalf("data source %q is registered more than once", registration.name)
		}
		dataSourceNames[registration.name] = struct{}{}

		response := &datasource.MetadataResponse{}
		registration.factory().Metadata(
			context.Background(),
			datasource.MetadataRequest{ProviderTypeName: "aci"},
			response,
		)
		if response.TypeName != registration.name {
			t.Fatalf("data source registration %q has metadata type name %q", registration.name, response.TypeName)
		}
	}

	functionNames := make(map[string]struct{}, len(functionRegistrations))
	for _, registration := range functionRegistrations {
		if registration.factory == nil {
			t.Fatalf("function %q has a nil factory", registration.name)
		}
		if _, exists := functionNames[registration.name]; exists {
			t.Fatalf("function %q is registered more than once", registration.name)
		}
		functionNames[registration.name] = struct{}{}

		response := &function.MetadataResponse{}
		registration.factory().Metadata(context.Background(), function.MetadataRequest{}, response)
		if response.Name != registration.name {
			t.Fatalf("function registration %q has metadata name %q", registration.name, response.Name)
		}
	}
}

func resourceFactoryNames(t *testing.T, factories []func() resource.Resource) []string {
	t.Helper()
	names := make([]string, len(factories))
	for index, factory := range factories {
		if factory == nil {
			t.Fatal("resource factory is nil")
		}
		response := &resource.MetadataResponse{}
		factory().Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "aci"}, response)
		names[index] = response.TypeName
	}
	sort.Strings(names)
	return names
}

func dataSourceFactoryNames(t *testing.T, factories []func() datasource.DataSource) []string {
	t.Helper()
	names := make([]string, len(factories))
	for index, factory := range factories {
		if factory == nil {
			t.Fatal("data source factory is nil")
		}
		response := &datasource.MetadataResponse{}
		factory().Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "aci"}, response)
		names[index] = response.TypeName
	}
	sort.Strings(names)
	return names
}

func functionFactoryNames(t *testing.T, factories []func() function.Function) []string {
	t.Helper()
	names := make([]string, len(factories))
	for index, factory := range factories {
		if factory == nil {
			t.Fatal("function factory is nil")
		}
		response := &function.MetadataResponse{}
		factory().Metadata(context.Background(), function.MetadataRequest{}, response)
		names[index] = response.Name
	}
	sort.Strings(names)
	return names
}
