package leasemanager_test

import (
	. "github.com/mandarjog/leasedispenser/leasemanager"
	"github.com/mandarjog/leasedispenser/leasemanager/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mem Lease Provider Test", func() {
	var (
		FakeProvider LeaseProvider
		leaseReq     ProviderLeaseRequest
		leaseID      string
		leaseinfo    LeaseInfo
		respInfo     ProviderLeaseInfo
		err          error
		skuID        string
	)
	BeforeEach(func() {
		FakeProvider = fakes.NewMemProvider()
		skuID = "SSSSSSSS"
		leaseID = "LLLLLLLLL"
		leaseReq = ProviderLeaseRequest{
			LeaseID:  leaseID,
			Duration: 1000000,
		}
		leaseinfo = LeaseInfo{
			ID: leaseID,
			Req: LeaseRequest{
				SKU: skuID,
			},
		}
	})
	Describe("Request", func() {
		BeforeEach(func() {
			respInfo, err = FakeProvider.Request(leaseReq)
		})
		It("Should accept a new lease without error and mark status pending", func() {
			Ω(err).Should(BeNil())
			Ω(respInfo.StatusCode).Should(Equal(LeaseStatusPending))
			Ω(respInfo.ProviderLeaseID).ShouldNot(Equal(""))
		})
		It("Should complete after being polled via Info()", func() {
			Eventually(func() string {
				pli, _ := FakeProvider.Info(respInfo.ProviderLeaseID)
				return pli.StatusCode
			}, 3, 1).Should(Equal(LeaseStatusActive))
			pli, _ := FakeProvider.Info(respInfo.ProviderLeaseID)
			Ω(pli.ProviderLeaseID).ShouldNot(Equal(""))
		})
	})
})
