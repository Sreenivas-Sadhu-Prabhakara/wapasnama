package backend

import (
	"encoding/json"
	"fmt"
)

// Claim is one discrepancy raised against a supplier: short quantity, damaged
// goods, wrong rate, or a missing scheme discount.
type Claim struct {
	Supplier string  `json:"supplier"`
	Item     string  `json:"item"`
	Kind     string  `json:"kind"`   // short | damaged | wrong-rate | missing-scheme
	Amount   float64 `json:"amount"` // rupees claimed
	Status   string  `json:"status"` // open | received
}

var validKinds = map[string]bool{"short": true, "damaged": true, "wrong-rate": true, "missing-scheme": true}

// Validate reports whether the Claim is well formed.
func (c Claim) Validate() error {
	if c.Supplier == "" {
		return fmt.Errorf("supplier is required")
	}
	if !validKinds[c.Kind] {
		return fmt.Errorf("kind must be short, damaged, wrong-rate or missing-scheme")
	}
	if c.Amount < 0 {
		return fmt.Errorf("amount cannot be negative")
	}
	if c.Status != "open" && c.Status != "received" {
		return fmt.Errorf("status must be open or received")
	}
	return nil
}

// Summary aggregates the claim ledger.
type Summary struct {
	OpenAmount     float64 `json:"openAmount"`
	OpenCount      int     `json:"openCount"`
	RecoveredAmount float64 `json:"recoveredAmount"`
	RecoveredCount int     `json:"recoveredCount"`
}

// Summarize computes open vs recovered totals from stored records. Each record's
// Headline is the claim amount and Label is its status.
func Summarize(records []Record) Summary {
	var s Summary
	for _, r := range records {
		switch r.Label {
		case "open":
			s.OpenAmount += r.Headline
			s.OpenCount++
		case "received":
			s.RecoveredAmount += r.Headline
			s.RecoveredCount++
		}
	}
	return s
}

// claimFromJSON decodes and validates a claim, returning its headline+label.
func claimFromJSON(raw []byte) (float64, string, error) {
	var c Claim
	if err := json.Unmarshal(raw, &c); err != nil {
		return 0, "", fmt.Errorf("invalid json")
	}
	if err := c.Validate(); err != nil {
		return 0, "", err
	}
	return c.Amount, c.Status, nil
}
