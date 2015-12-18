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
		FakeDB        fakes.MemLeaseDB
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
	Describe("Request Happy path", func() {
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
		AfterEach(func() {
			FakeDB.Clear()
		})
	})
	Describe("Given a successful lease request for 2 skus", func() {
		BeforeEach(func() {
			leaseManager.Request(req)
			leaseManager.Request(req)
		})
		It("Then 2 leases Should be obtained thru DB.FindAll()", func() {
			leases, err := FakeDB.FindAll()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(leases)).Should(Equal(2))
		})
		It("Then 2 leases should be obtained thru LeaseManager.List() call", func() {
			leases, err := leaseManager.List("", "")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(leases)).Should(Equal(2))
		})
		AfterEach(func() {
			FakeDB.Clear()
		})
	})
	Describe("Request for bad SKU", func() {
		BeforeEach(func() {
			respLeaseinfo, err = leaseManager.Request(LeaseRequest{SKU: "DOESNOTEXIST"})
		})
		It("Should Fail with error message", func() {
			Ω(err).ShouldNot(BeNil())
			Ω(err.Error()).Should(ContainSubstring("Not found"))
		})
	})
})
