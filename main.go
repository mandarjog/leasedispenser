package main

import (
	. "github.com/mandarjog/leasedispenser/leasemanager"
	"github.com/mandarjog/leasedispenser/leasemanager/fakes"
)

func main() {

}

// BuildLeaseManager -- returns a fully constructed
// LeaseManager This can be driven by a CLI or
// a web App
func BuildLeaseManager() LeaseManager {
	return NewLeaseManager(BuildProviderRegistry(), BuildLeaseDB())
}

// BuildProviderRegistry -- instantiate All providers and
// construct a registry
func BuildProviderRegistry() map[string]LeaseProvider {
	registry := make(map[string]LeaseProvider)
	registry["fakeSKU"] = fakes.NewMemProvider()
	return registry
}

// BuildLeaseDB -- create a highlevel database provider
func BuildLeaseDB() LeaseDB {
	return fakes.NewLeaseDB()
}

// BuildRouter -- create a router based on above
