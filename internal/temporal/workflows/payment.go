package workflows

import (
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
)

type PaymentDetails struct {
	PayID  uuid.UUID
	Amount float64
}

func ProcessPayment(ctx workflow.Context, in PaymentDetails) error {
	// TODO: add payment processing logic.
	workflow.GetLogger(ctx).Info("Did some processing and recieved the payment")

	return nil
}
