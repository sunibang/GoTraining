package activities

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/romangurevitch/go-training/internal/temporal/order"
	"go.temporal.io/sdk/temporal"
)

type OrderActivities struct {
	inventoryClient InventoryChecker
}

func NewOrderActivities(inventoryClient InventoryChecker) *OrderActivities {
	return &OrderActivities{
		inventoryClient: inventoryClient,
	}
}

func (a *OrderActivities) Validate(ctx context.Context, o order.Order) error {
	if (o.ID == uuid.UUID{}) {
		return fmt.Errorf("order must have a valid order ID")
	}

	if len(o.LineItems) < 1 {
		return fmt.Errorf("order must have at least one item")
	}

	// Check inventory for each line item
	for _, item := range o.LineItems {
		available, err := a.inventoryClient.CheckInventory(ctx, item.ProductID, item.Quantity)
		if err != nil {
			return fmt.Errorf("failed to check inventory for product %s: %w", item.ProductID, err)
		}
		if !available {
			return temporal.NewNonRetryableApplicationError(
				"insufficient inventory for product",
				"validation",
				fmt.Errorf("insufficient inventory for product %s", item.ProductID),
			)
		}
	}

	return nil
}

func (a *OrderActivities) Process(ctx context.Context, o order.Order) (string, error) {
	// TODO: add order processing logic.

	return "Processed", nil
}

func (a *OrderActivities) Pick(ctx context.Context, o order.Order) error {
	// Simulate warehouse picking.
	slog.InfoContext(ctx, "Picking order", "order_id", o.ID)
	return nil
}

func (a *OrderActivities) Ship(ctx context.Context, o order.Order) error {
	// Simulate creating a shipment.
	slog.InfoContext(ctx, "Shipping order", "order_id", o.ID)
	return nil
}

func (a *OrderActivities) Deliver(ctx context.Context, o order.Order) error {
	// Simulate delivery confirmation.
	slog.InfoContext(ctx, "Delivering order", "order_id", o.ID)
	return nil
}
