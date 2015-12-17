package fakes

import (
	"errors"
	"time"

	"code.google.com/p/go-uuid/uuid"
	. "github.com/mandarjog/leasedispenser/leasemanager"
)

type MemProvider struct {
	DB map[string]ProviderLeaseInfo
}

func NewMemProvider() LeaseProvider {
	return MemProvider{
		DB: make(map[string]ProviderLeaseInfo),
	}
}

func (s MemProvider) Request(req ProviderLeaseRequest) (info ProviderLeaseInfo, err error) {
	providerLeaseID := uuid.New()
	info.ProviderLeaseID = providerLeaseID
	info.StatusCode = LeaseStatusPending
	newinfo := info
	newinfo.StatusCode = LeaseStatusActive
	newinfo.LeaseStartDate = time.Now().UnixNano()
	newinfo.LeaseEndDate = newinfo.LeaseStartDate + req.Duration
	s.DB[providerLeaseID] = newinfo
	return
}
func (s MemProvider) Info(providerLeaseID string) (info ProviderLeaseInfo, err error) {
	var found bool
	if info, found = s.DB[providerLeaseID]; !found {
		err = errors.New(providerLeaseID + " was not found")
	}
	return
}
