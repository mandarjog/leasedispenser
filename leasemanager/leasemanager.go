package leasemanager

import (
	"errors"

	"code.google.com/p/go-uuid/uuid"
)

type (
	// LeaseManagerState -- holds data for the LeaseManager
	LeaseManagerState struct {
		Registry map[string]LeaseProvider
		DB       LeaseDB
	}
)

// NewLeaseManager -- construct a new lease manager with given params
func NewLeaseManager(registry map[string]LeaseProvider, db LeaseDB) LeaseManager {
	s := LeaseManagerState{
		Registry: registry,
		DB:       db,
	}
	return s
}

// Info -- provide info using providers
// If usecache is false, it should go back to the provider to ask for an update
// else it can return results from the DB
func (s LeaseManagerState) Info(leaseID string, usecache bool) (LeaseInfo, error) {
	var (
		found    bool
		provider LeaseProvider
		err      error
		info     LeaseInfo
		pli      ProviderLeaseInfo
		nullinfo = LeaseInfo{}
	)

	if info, err = s.DB.FindByID(leaseID); err != nil || usecache {
		return info, err
	}

	if provider, found = s.Registry[info.Req.SKU]; !found {
		return nullinfo, errors.New(info.Req.SKU + " Not found")
	}

	if pli, err = provider.Info(info.ProviderLeaseID); err != nil {
		return nullinfo, err
	}

	leaseInfo := LeaseInfo{
		pli,
		leaseID,
		info.Req,
	}
	err = s.DB.CreateOrUpdate(leaseID, leaseInfo)
	return leaseInfo, err
}

// List -- List all leases with 2 selectors
// empty selector matches everything
func (s LeaseManagerState) List(sku string, owner string) (info []LeaseInfo, err error) {
	leases, err := s.DB.FindAll()
	if err == nil {
		for _, lease := range leases {
			if (sku == "" || lease.Req.SKU == sku) && (owner == "" || lease.Req.Owner == owner) {
				info = append(info, lease)
			}
		}
	}
	return
}

func (s LeaseManagerState) PollForPendingLeases(pollInterval int64, maxPolls int) {
	for idx := 0; idx <= maxPolls; idx++ {
	}

}

// Request -- use providers to satisfy requests
func (s LeaseManagerState) Request(req LeaseRequest) (LeaseInfo, error) {
	var (
		found        bool
		provider     LeaseProvider
		leaseID      string
		err          error
		providerInfo ProviderLeaseInfo
		nullinfo     = LeaseInfo{}
	)

	if provider, found = s.Registry[req.SKU]; !found {
		return nullinfo, errors.New(req.SKU + " Not found")
	}

	//FIXME check for request duplication
	leaseID = uuid.New()
	leaseInfo := LeaseInfo{
		ID:  leaseID,
		Req: req,
	}
	leaseInfo.StatusCode = LeaseStatusRequested
	if err = s.DB.CreateOrUpdate(leaseID, leaseInfo); err != nil {
		return nullinfo, err
	}

	pReq := ProviderLeaseRequest{
		LeaseID:  leaseID,
		Owner:    req.Owner,
		Duration: req.Duration,
		ReqInfo:  req.ReqInfo,
	}

	if providerInfo, err = provider.Request(pReq); err != nil {
		return leaseInfo, err
	}

	leaseInfo = LeaseInfo{
		providerInfo,
		leaseID,
		req,
	}

	err = s.DB.CreateOrUpdate(leaseID, leaseInfo)
	return leaseInfo, err
}
