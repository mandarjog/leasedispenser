package leasemanager

import (
	"errors"

	"code.google.com/p/go-uuid/uuid"
	"github.com/xchapter7x/lo"
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
	return LeaseManagerState{
		registry,
		db,
	}
}

// Info -- provide info using providers
// If usecache is false, it should go back to the provider to ask for an update
// else it can return results from the DB
func (s LeaseManagerState) Info(leaseID string, usecache bool) (LeaseInfo, error) {
	var (
		found    bool
		provider LeaseProvider
		err      error
	)
	info, err := s.DB.FindByID(leaseID)
	if err != nil {
		return info, err
	}

	if !usecache {
		provider, found = s.Registry[info.Req.SKU]
		if !found {
			err = errors.New(info.Req.SKU + " Not found")
			lo.G.Error(err.Error())
			return LeaseInfo{}, err
		}

		pli, err := provider.Info(info.ProviderLeaseID)
		if err != nil {
			lo.G.Error(err.Error())
			return LeaseInfo{}, err
		}
		leaseInfo := LeaseInfo{
			pli,
			leaseID,
			info.Req,
		}

		err = s.DB.CreateOrUpdate(leaseID, leaseInfo)
		if err != nil {
			lo.G.Error(err.Error())
		}

		info = leaseInfo

	}
	return info, err
}

// Request -- use providers to satisfy requests
func (s LeaseManagerState) Request(req LeaseRequest) (LeaseInfo, error) {
	var (
		found    bool
		provider LeaseProvider
		leaseID  string
		err      error
	)
	provider, found = s.Registry[req.SKU]
	if !found {
		err = errors.New(req.SKU + " Not found")
		lo.G.Error(err.Error())
		return LeaseInfo{}, err
	}
	//TODO check for request duplication
	leaseID = uuid.New()
	leaseInfo := LeaseInfo{
		ID:  leaseID,
		Req: req,
	}
	leaseInfo.StatusCode = LeaseStatusRequested
	err = s.DB.CreateOrUpdate(leaseID, leaseInfo)
	if err != nil {
		lo.G.Error(err.Error())
		return LeaseInfo{}, err
	}

	pReq := ProviderLeaseRequest{
		LeaseID:  leaseID,
		Owner:    req.Owner,
		Duration: req.Duration,
		ReqInfo:  req.ReqInfo,
	}

	providerInfo, err := provider.Request(pReq)
	if err != nil {
		lo.G.Error(err.Error())
		return leaseInfo, err
	}
	leaseInfo = LeaseInfo{
		providerInfo,
		leaseID,
		req,
	}

	err = s.DB.CreateOrUpdate(leaseID, leaseInfo)
	if err != nil {
		lo.G.Error(err.Error())
	}

	return leaseInfo, err
}
