package leasemanager_test

import (
	. "github.com/mandarjog/leasedispenser/leasemanager"
	"github.com/mandarjog/leasedispenser/leasemanager/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MemDB Provider Test", func() {
	var (
		FakeDB    LeaseDB
		leaseID   string
		leaseinfo LeaseInfo
		err       error
		skuID     string
	)
	BeforeEach(func() {
		FakeDB = fakes.NewLeaseDB()
		skuID = "SSSSSSSS"
		leaseID = "LLLLLLLLL"
		leaseinfo = LeaseInfo{
			ID: leaseID,
			Req: LeaseRequest{
				SKU: skuID,
			},
		}
	})
	Describe("CreateorUpdate", func() {
		BeforeEach(func() {
			err = FakeDB.CreateOrUpdate(leaseID, leaseinfo)
		})
		It("Should accept a new lease and leasinfo without error", func() {
			Ω(err).Should(BeNil())
		})
		It("Should find the lease that was just entered", func() {
			foundlease, err1 := FakeDB.FindByID(leaseID)
			Ω(err1).Should(BeNil())
			Ω(foundlease).Should(Equal(leaseinfo))
		})
		It("Should NOT find the lease that was not entered", func() {
			foundlease, err1 := FakeDB.FindByID("BOBOBO")
			Ω(err1).ShouldNot(BeNil())
			Ω(foundlease).ShouldNot(Equal(leaseinfo))
		})
	})
	Describe("FindBySKU", func() {
		BeforeEach(func() {
			FakeDB.CreateOrUpdate(leaseID, leaseinfo)
		})
		It("Should  find by valid sku", func() {
			leases, err1 := FakeDB.FindBySKU(skuID)
			Ω(err1).Should(BeNil())
			Ω(len(leases)).Should(Equal(1))
		})
		It("Should NOT find the lease that was not entered", func() {
			leases, err1 := FakeDB.FindBySKU("DOESNOTEXIST")
			Ω(err1).Should(BeNil())
			Ω(len(leases)).Should(Equal(0))
		})
	})
})
