package customers

import "testing"

func TestEstimatedAnnualRevenueTopBandMatchesWire(t *testing.T) {
	// The band below ends at 249999999, so the top band is 250 million, and the
	// spec enum on CreateCustomerIn agrees.
	if got := string(EstimatedAnnualRevenue250000000Plus); got != "250000000_plus" {
		t.Fatalf("top band = %q, want 250000000_plus", got)
	}
	if got := string(EstimatedAnnualRevenue2500000000Plus); got != "250000000_plus" {
		t.Fatalf("deprecated alias = %q, want the corrected wire value", got)
	}
}
