package fakes

import (
	"errors"

	. "github.com/mandarjog/leasedispenser/leasemanager"
)

type MemLeaseDB struct {
	DB map[string]LeaseInfo
}

func NewLeaseDB() MemLeaseDB {
	return MemLeaseDB{
		DB: make(map[string]LeaseInfo),
	}
}

func (s MemLeaseDB) Clear() {
	s.DB = make(map[string]LeaseInfo)
}

func (s MemLeaseDB) CreateOrUpdate(leaseID string, lease LeaseInfo) (err error) {
	s.DB[leaseID] = lease
	return
}

func (s MemLeaseDB) FindByID(leaseID string) (lease LeaseInfo, err error) {
	lease, found := s.DB[leaseID]
	if !found {
		err = errors.New("Unable to find " + leaseID)
	}
	return
}

func (s MemLeaseDB) FindBySKU(sku string) (leases []LeaseInfo, err error) {
	for _, lease := range s.DB {
		if lease.Req.SKU == sku {
			leases = append(leases, lease)
		}
	}
	return
}

func (s MemLeaseDB) FindAll() (leases []LeaseInfo, err error) {
	for _, lease := range s.DB {
		leases = append(leases, lease)
	}
	return
}
