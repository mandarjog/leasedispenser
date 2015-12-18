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
		Request(req LeaseRequest) (info LeaseInfo, err error)
		Info(leaseID string, usecache bool) (info LeaseInfo, err error)
		List(sku string, owner string) (info []LeaseInfo, err error)
		//PollOutstandingLeases()
	}

	LeaseDB interface {
		CreateOrUpdate(leaseID string, lease LeaseInfo) (err error)
		FindBySKU(sku string) (leases []LeaseInfo, err error)
		FindByID(leaseID string) (lease LeaseInfo, err error)
		FindAll() (leases []LeaseInfo, err error)
	}
)
