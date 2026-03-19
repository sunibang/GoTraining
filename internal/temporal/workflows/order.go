package workflows

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/romangurevitch/go-training/internal/temporal/activities"
	"github.com/romangurevitch/go-training/internal/temporal/order"
	"go.temporal.io/sdk/workflow"
)

type Params struct {
	Order order.Order
}

func ProcessOrder(ctx workflow.Context, in Params) (order.OrderStatus, error) {
	logger := workflow.GetLogger(ctx)

	var orderStatus order.OrderStatus

	err := workflow.SetQueryHandler(ctx, "GetOrderStatus", func() (order.OrderStatus, error) {
		return orderStatus, nil
	})
	if err != nil {
		return orderStatus, fmt.Errorf("failed to setup query handler: %w", err)
	}

	// Validate order and items.
	ctx = workflow.WithActivityOptions(ctx, validateActivityOptions)

	var orderActivities *activities.OrderActivities
	err = workflow.ExecuteActivity(ctx, orderActivities.Validate, in.Order).Get(ctx, nil)
	if err != nil {
		orderStatus = order.UnableToComplete
		return orderStatus, err
	}
	orderStatus = order.Placed

	// Wait for order picked or order cancelled signals.
	pickOrderCh := workflow.GetSignalChannel(ctx, pickOrderSignal)
	cancelOrderCh := workflow.GetSignalChannel(ctx, cancelOrderSignal)

	selector := workflow.NewSelector(ctx)
	selector.AddReceive(pickOrderCh, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, nil)
		now := workflow.Now(ctx)
		logger.Info("Order picked at %s", now.Format("2006-01-02 15:04:05"))
		orderStatus = order.Picked
	})
	selector.AddReceive(cancelOrderCh, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, nil)
		orderStatus = order.Cancelled
	})

	// Blocks until signal is received.
	selector.Select(ctx)
	if orderStatus == order.Cancelled {
		workflow.GetLogger(ctx).Warn("Received cancellation signal")
		return orderStatus, nil
	}

	// Process order.
	var payID uuid.UUID
	if err = workflow.SideEffect(ctx, func(ctx workflow.Context) any {
		return uuid.New()
	}).Get(&payID); err != nil {
		orderStatus = order.UnableToComplete
		return orderStatus, err
	}

	if err = workflow.ExecuteChildWorkflow(ctx, ProcessPayment, PaymentDetails{
		PayID:  payID,
		Amount: 100.0, // TODO: calculate amount.
	}).Get(ctx, nil); err != nil {
		orderStatus = order.UnableToComplete
		return orderStatus, err
	}

	// Wait for order to be shipped.
	workflow.GetSignalChannel(ctx, shipOrderSignal).Receive(ctx, nil)
	orderStatus = order.Shipped

	// Wait for order to be marked as delivered.
	workflow.GetSignalChannel(ctx, orderDeliveredSignal).Receive(ctx, nil)
	orderStatus = order.Completed

	return orderStatus, nil
}

// AutoProcessOrder is a simplified version of ProcessOrder that drives the
// order through each lifecycle stage automatically via activities, removing
// the need for external signals.
//
// Lifecycle: PLACED → PICKED → SHIPPED → COMPLETED
func AutoProcessOrder(ctx workflow.Context, in Params) (order.OrderStatus, error) {
	// status tracks the current stage of the order. It is updated as each
	// stage completes, and is exposed via the query handler below.
	var status order.OrderStatus

	// Register a query handler so external clients can ask "where is my order?"
	// at any point during execution. Temporal replays the workflow history to
	// answer queries, so this must be set up before any blocking calls.
	if err := workflow.SetQueryHandler(ctx, "GetOrderStatus", func() (order.OrderStatus, error) {
		return status, nil
	}); err != nil {
		return status, fmt.Errorf("failed to register query handler: %w", err)
	}

	// a is a typed nil used only to reference activity methods by name.
	// Temporal resolves the real implementation (registered in the worker)
	// from this function reference — it never calls through the nil pointer.
	var a *activities.OrderActivities

	// --- Stage 1: Validate ---
	// Use the tighter timeout and higher attempt count suited for fast,
	// external validation calls.
	validateCtx := workflow.WithActivityOptions(ctx, validateActivityOptions)
	if err := workflow.ExecuteActivity(validateCtx, a.Validate, in.Order).Get(ctx, nil); err != nil {
		status = order.UnableToComplete
		return status, fmt.Errorf("validation failed: %w", err)
	}
	status = order.Placed

	// --- Stage 2: Pick ---
	// Apply the default (longer) timeout for warehouse operations.
	defaultCtx := workflow.WithActivityOptions(ctx, defaultActivityOptions)
	if err := workflow.ExecuteActivity(defaultCtx, a.Pick, in.Order).Get(ctx, nil); err != nil {
		status = order.UnableToComplete
		return status, fmt.Errorf("picking failed: %w", err)
	}
	status = order.Picked

	// --- Stage 3: Payment ---
	// Generate a unique payment ID using SideEffect so that the ID is stable
	// across workflow replays. Without SideEffect, uuid.New() would produce a
	// different value each time the workflow history is replayed, breaking
	// Temporal's determinism requirement.
	var payID uuid.UUID
	if err := workflow.SideEffect(ctx, func(_ workflow.Context) any {
		return uuid.New()
	}).Get(&payID); err != nil {
		status = order.UnableToComplete
		return status, fmt.Errorf("failed to generate payment ID: %w", err)
	}

	// ProcessPayment runs as a child workflow so its execution history is
	// tracked separately. This makes it independently observable and retryable.
	if err := workflow.ExecuteChildWorkflow(ctx, ProcessPayment, PaymentDetails{
		PayID:  payID,
		Amount: 100.0, // TODO: calculate from order line items.
	}).Get(ctx, nil); err != nil {
		status = order.UnableToComplete
		return status, fmt.Errorf("payment failed: %w", err)
	}

	// --- Stage 4: Ship ---
	if err := workflow.ExecuteActivity(defaultCtx, a.Ship, in.Order).Get(ctx, nil); err != nil {
		status = order.UnableToComplete
		return status, fmt.Errorf("shipping failed: %w", err)
	}
	status = order.Shipped

	// --- Stage 5: Deliver ---
	if err := workflow.ExecuteActivity(defaultCtx, a.Deliver, in.Order).Get(ctx, nil); err != nil {
		status = order.UnableToComplete
		return status, fmt.Errorf("delivery failed: %w", err)
	}
	status = order.Completed

	return status, nil
}
