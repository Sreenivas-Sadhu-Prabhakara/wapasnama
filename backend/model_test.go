package backend

import (
	"encoding/json"
	"testing"
)

func rec(amount float64, status string) Record {
	in, _ := json.Marshal(Claim{Supplier: "ACME", Item: "x", Kind: "short", Amount: amount, Status: status})
	return Record{Input: in, Headline: amount, Label: status}
}

func TestSummarize(t *testing.T) {
	s := Summarize([]Record{rec(100, "open"), rec(250, "open"), rec(80, "received")})
	if s.OpenCount != 2 || s.OpenAmount != 350 {
		t.Fatalf("open agg wrong: %+v", s)
	}
	if s.RecoveredCount != 1 || s.RecoveredAmount != 80 {
		t.Fatalf("recovered agg wrong: %+v", s)
	}
}

func TestClaimValidate(t *testing.T) {
	ok := Claim{Supplier: "ACME", Kind: "short", Amount: 10, Status: "open"}
	if err := ok.Validate(); err != nil {
		t.Fatalf("valid rejected: %v", err)
	}
	for i, bad := range []Claim{
		{Kind: "short", Status: "open"},                       // no supplier
		{Supplier: "A", Kind: "nope", Status: "open"},         // bad kind
		{Supplier: "A", Kind: "short", Status: "maybe"},       // bad status
		{Supplier: "A", Kind: "short", Amount: -1, Status: "open"},
	} {
		if err := bad.Validate(); err == nil {
			t.Fatalf("bad claim %d accepted", i)
		}
	}
}
