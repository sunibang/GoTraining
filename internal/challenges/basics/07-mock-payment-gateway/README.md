# Challenge 07: Mock a Payment Gateway

**Difficulty:** Intermediate
**Covers:** `mocking`, `interface`, `testing`

---

## Goal

Build a `PaymentService` that depends on a payment gateway through a consumer-side interface. Use GoMock to write tests that verify call sequencing and error handling.

---

## Requirements

### 1. Define the interface (consumer-side)

In the `payment` package (the consumer), define:

```go
// Gateway is defined here, in the consumer package — not in the gateway package.
type Gateway interface {
    Charge(amount int) (string, error) // returns transaction ID
    Refund(id string) error
}
```

### 2. Implement `PaymentService`

```go
type PaymentService struct {
    gateway Gateway
    charged map[string]bool // tracks charged transaction IDs
}

func NewPaymentService(g Gateway) *PaymentService

// Charge charges the gateway and records the transaction ID.
func (s *PaymentService) Charge(amount int) (string, error)

// Refund refunds a previously charged transaction.
// Returns an error if the transaction ID was never charged (double-refund prevention).
func (s *PaymentService) Refund(id string) error
```

### 3. Generate a mock

```bash
mockgen -destination=mocks/mock_gateway.go -package=mocks \
    github.com/romangurevitch/go-training/internal/challenges/basics/07-mock-payment-gateway/payment Gateway
```

### 4. Write tests for

| Test case | GoMock feature |
|-----------|---------------|
| Successful charge | `EXPECT().Charge(100).Return("tx-1", nil).Times(1)` |
| Gateway returns error on charge | `Return("", errors.New("declined"))` |
| Successful refund after charge | `After(chargeCall)` to enforce order |
| Double-refund prevention (no gateway call) | Assert `Refund` is **never** called on mock |

---

## Skills Practiced

- Consumer-side interface definition
- GoMock: `Times(n)`, `After(call)`, `AnyTimes()`, `Return()`
- Test isolation: each test creates its own controller and mock
- Preventing double-side-effects in business logic

---

## Hints

- `gomock.NewController(t)` — create one per test function, not per package
- `EXPECT().Method().Times(0)` — asserts a method is **never** called
- Use a `map[string]bool` to track charged IDs in `PaymentService`
- For the double-refund test, the mock's `Refund` should have `Times(0)` — the service should return an error before ever calling the gateway
