package leasemanager_test

import (
	. "github.com/mandarjog/leasedispenser/leasemanager"
	"github.com/mandarjog/leasedispenser/leasemanager/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lease Manager Test", func() {
	var (
		leaseManager  LeaseManager
		FakeProvider  LeaseProvider
		FakeDB        LeaseDB
		leaseReq      ProviderLeaseRequest
		leaseID       string
		leaseinfo     LeaseInfo
		respLeaseinfo LeaseInfo
		err           error
		skuID         string
		registry      map[string]LeaseProvider
		req           LeaseRequest
	)
	BeforeEach(func() {
		FakeProvider = fakes.NewMemProvider()
		FakeDB = fakes.NewLeaseDB()
		skuID = "SSSSSSSS"
		leaseID = "LLLLLLLLL"
		leaseReq = ProviderLeaseRequest{
			LeaseID:  leaseID,
			Duration: 1000000,
		}
		req = LeaseRequest{
			SKU: skuID,
		}
		leaseinfo = LeaseInfo{
			ID:  leaseID,
			Req: req,
		}
		registry = make(map[string]LeaseProvider)
		registry[skuID] = FakeProvider
		leaseManager = NewLeaseManager(registry, FakeDB)
	})
	Describe("Request", func() {
		BeforeEach(func() {
			respLeaseinfo, err = leaseManager.Request(req)
		})
		It("Should accept a new lease without error and mark status pending", func() {
			Ω(err).Should(BeNil())
			Ω(respLeaseinfo.StatusCode).Should(Equal(LeaseStatusPending))
			Ω(respLeaseinfo.ProviderLeaseID).ShouldNot(Equal(""))
		})
		It("Should complete after being polled via Info()", func() {
			Eventually(func() string {
				pli, _ := leaseManager.Info(respLeaseinfo.ID, false)
				return pli.StatusCode
			}, 3, 1).Should(Equal(LeaseStatusActive))
			pli, _ := leaseManager.Info(respLeaseinfo.ID, false)
			Ω(pli.ProviderLeaseID).ShouldNot(Equal(""))
		})
	})
})
