package activities_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/romangurevitch/go-training/internal/temporal/activities"
	"github.com/romangurevitch/go-training/internal/temporal/activities/mocks"
	"github.com/romangurevitch/go-training/internal/temporal/order"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
	"go.uber.org/mock/gomock"
)

const dummyOrderID = "8c727b70-cfcb-4674-8bcd-78e66e32f723"

func TestActivities(t *testing.T) {
	suite.Run(t, new(ActivityTestSuite))
}

type ActivityTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestActivityEnvironment
}

func (s *ActivityTestSuite) SetupTest() {
	s.env = s.NewTestActivityEnvironment()
}

func (s *ActivityTestSuite) TestValidate_Success() {
	// Setup
	inventoryChecker := mocks.NewMockInventoryChecker(gomock.NewController(s.T()))
	inventoryChecker.EXPECT().
		CheckInventory(gomock.Any(), uuid.MustParse("ba320a5d-62ed-46d0-b491-084514598721"), int32(1)).
		Return(true, nil)

	acts := activities.NewOrderActivities(inventoryChecker)
	s.env.RegisterActivity(acts.Validate)

	// Invoke
	_, err := s.env.ExecuteActivity(acts.Validate, order.Order{
		ID: uuid.MustParse(dummyOrderID),
		LineItems: []order.LineItem{
			{
				ProductID:    uuid.MustParse("ba320a5d-62ed-46d0-b491-084514598721"),
				Quantity:     1,
				PricePerItem: decimal.RequireFromString("123.45"),
			},
		},
	})

	// Assert
	s.Require().NoError(err)
}

func (s *ActivityTestSuite) TestValidate_Fail() {

	tests := []struct {
		name       string
		input      order.Order
		setupMocks func(t *testing.T, mockIC *mocks.MockInventoryChecker)
		err        string
	}{
		{
			name:  "Missing order ID",
			input: order.Order{},
			err:   "order must have a valid order ID",
		},
		{
			name: "No items in order",
			input: order.Order{
				ID:        uuid.MustParse(dummyOrderID),
				LineItems: []order.LineItem{},
			},
			err: "order must have at least one item",
		},
		{
			name: "No inventory for line item",
			input: order.Order{
				ID: uuid.MustParse(dummyOrderID),
				LineItems: []order.LineItem{
					{
						ProductID:    uuid.MustParse("ba320a5d-62ed-46d0-b491-084514598721"),
						Quantity:     1,
						PricePerItem: decimal.RequireFromString("123.45"),
					},
				},
			},
			setupMocks: func(t *testing.T, mockIC *mocks.MockInventoryChecker) {
				mockIC.EXPECT().CheckInventory(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
			},
			err: "insufficient inventory for product",
		},
		{
			name: "Inventory checker error",
			input: order.Order{
				ID: uuid.MustParse(dummyOrderID),
				LineItems: []order.LineItem{
					{
						ProductID:    uuid.MustParse("ba320a5d-62ed-46d0-b491-084514598721"),
						Quantity:     1,
						PricePerItem: decimal.RequireFromString("123.45"),
					},
				},
			},
			setupMocks: func(t *testing.T, mockIC *mocks.MockInventoryChecker) {
				mockIC.EXPECT().CheckInventory(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, errors.New("test error"))
			},
			err: "failed to check inventory for product",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Setup
			inventoryChecker := mocks.NewMockInventoryChecker(gomock.NewController(s.T()))
			if tt.setupMocks != nil {
				tt.setupMocks(s.T(), inventoryChecker)
			}

			acts := activities.NewOrderActivities(inventoryChecker)
			s.env.RegisterActivity(acts.Validate)

			// Invoke
			_, err := s.env.ExecuteActivity(acts.Validate, tt.input)

			// Assert
			s.Require().ErrorContains(err, tt.err)
		})
	}
}
