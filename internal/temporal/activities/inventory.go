package activities

import (
	"context"

	"github.com/google/uuid"
)

type InventoryChecker interface {
	CheckInventory(context.Context, uuid.UUID, int32) (bool, error)
}
