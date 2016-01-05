package leasemanager

type (
	ProviderLeaseInfo struct {
		ProviderLeaseID string
		StatusCode      string
		StatusMessage   string
		Info            map[string]string
		LeaseEndDate    int64
		LeaseStartDate  int64
	}
	ProviderLeaseRequest struct {
		LeaseID  string
		Owner    string
		Duration int64
		ReqInfo  map[string]interface{}
	}
	LeaseProvider interface {
		Request(req ProviderLeaseRequest) (info ProviderLeaseInfo, err error)
		Info(providerLeaseID string) (info ProviderLeaseInfo, err error)
		/*
			Release
			Update
			List
		*/
	}
	LeaseRequest struct {
		SKU      string
		Owner    string
		Duration int64
		ReqInfo  map[string]interface{}
	}
	LeaseInfo struct {
		ProviderLeaseInfo
		ID  string
		Req LeaseRequest
	}

	LeaseManager interface {
		// NewLeaseManager(registry map[string]LeaseProvider, db LeaseDB) (leaseManager LeaseManager)
		// Request -- Request a Lease and schedule call to the back end
		Request(req LeaseRequest) (info LeaseInfo, err error)
		// Info -- Get info about the named lease. If usecache=false, attempt to contact  provider for pending leases
		Info(leaseID string, usecache bool) (info LeaseInfo, err error)
		// List -- List lease info for leases matching optional sku and owner
		List(sku string, owner string) (info []LeaseInfo, err error)
		// PollForPendingLeases -- Attempt to contact providers for pending leases and update database
		// It should poll at most maxPolls times before exiting
		PollForPendingLeases(pollInterval int64, maxPolls int)
	}

	LeaseDB interface {
		CreateOrUpdate(leaseID string, lease LeaseInfo) (err error)
		FindBySKU(sku string) (leases []LeaseInfo, err error)
		FindByID(leaseID string) (lease LeaseInfo, err error)
		FindAll() (leases []LeaseInfo, err error)
	}
)
