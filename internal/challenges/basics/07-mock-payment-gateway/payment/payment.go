// Package payment is a challenge skeleton. Complete the TODOs.
package payment

// Gateway is the consumer-side interface for charging and refunding.
// Defined here — in the consumer package — so tests can mock it locally.
type Gateway interface {
	Charge(amount int) (string, error) // returns transaction ID
	Refund(id string) error
}

// PaymentService orchestrates payments through a Gateway.
type PaymentService struct {
	gateway Gateway
	charged map[string]bool // tracks IDs that have been charged
}

// NewPaymentService returns a PaymentService backed by the given Gateway.
func NewPaymentService(g Gateway) *PaymentService {
	return &PaymentService{
		gateway: g,
		charged: make(map[string]bool),
	}
}

// Charge delegates to the gateway and records the transaction ID on success.
func (s *PaymentService) Charge(amount int) (string, error) {
	// TODO: call s.gateway.Charge(amount), store txID in s.charged on success
	panic("not implemented")
}

// Refund refunds a previously charged transaction.
// Returns an error if the ID was never charged (prevents double-refund).
func (s *PaymentService) Refund(id string) error {
	// TODO:
	// 1. Check s.charged[id] — return an error if not found
	// 2. Call s.gateway.Refund(id)
	// 3. Remove id from s.charged on success
	panic("not implemented")
}
