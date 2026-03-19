package workflows_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/romangurevitch/go-training/internal/temporal/activities"
	"github.com/romangurevitch/go-training/internal/temporal/order"
	"github.com/romangurevitch/go-training/internal/temporal/workflows"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	suite.Run(t, new(WorkflowTestSuite))
}

type WorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment

	activities *activities.OrderActivities
}

func (s *WorkflowTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
	// ProcessPayment is used as a child workflow and must be explicitly
	// registered so the test environment can resolve it during execution.
	s.env.RegisterWorkflow(workflows.ProcessPayment)
	s.activities = &activities.OrderActivities{}
}

func (s *WorkflowTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *WorkflowTestSuite) TestWorkflow_Success() {
	// Mock activity implementations.

	s.env.OnActivity(s.activities.Validate, mock.Anything, order.Order{}).Return(nil)

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow("pickOrder", nil)
	}, time.Minute)

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow("shipOrder", nil)
	}, 2*time.Hour)
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow("markOrderAsDelivered", nil)
	}, 5*24*time.Hour)

	// Execute workflow.

	s.env.ExecuteWorkflow(workflows.ProcessOrder, workflows.Params{Order: order.Order{}})

	// Assert execution and order status.

	s.Require().NoError(s.env.GetWorkflowError())

	val, err := s.env.QueryWorkflow("GetOrderStatus")
	s.Require().NoError(err, "workflow should be queryable")
	var got order.OrderStatus
	err = val.Get(&got)
	s.Require().NoError(err, "query result should be an order.OrderStatus")
	s.Equal(order.Completed, got, "order should be completed")
}

func (s *WorkflowTestSuite) TestAutoWorkflow_Success() {
	s.env.OnActivity(s.activities.Validate, mock.Anything, order.Order{}).Return(nil)
	s.env.OnActivity(s.activities.Pick, mock.Anything, order.Order{}).Return(nil)
	s.env.OnActivity(s.activities.Ship, mock.Anything, order.Order{}).Return(nil)
	s.env.OnActivity(s.activities.Deliver, mock.Anything, order.Order{}).Return(nil)

	s.env.ExecuteWorkflow(workflows.AutoProcessOrder, workflows.Params{Order: order.Order{}})

	s.Require().NoError(s.env.GetWorkflowError())

	val, err := s.env.QueryWorkflow("GetOrderStatus")
	s.Require().NoError(err)
	var got order.OrderStatus
	s.Require().NoError(val.Get(&got))
	s.Equal(order.Completed, got)
}

func (s *WorkflowTestSuite) TestAutoWorkflow_ValidationFailure() {
	validationErr := fmt.Errorf("invalid order")
	s.env.OnActivity(s.activities.Validate, mock.Anything, order.Order{}).Return(validationErr)

	s.env.ExecuteWorkflow(workflows.AutoProcessOrder, workflows.Params{Order: order.Order{}})

	s.Require().Error(s.env.GetWorkflowError())

	val, err := s.env.QueryWorkflow("GetOrderStatus")
	s.Require().NoError(err)
	var got order.OrderStatus
	s.Require().NoError(val.Get(&got))
	s.Equal(order.UnableToComplete, got)
}

func (s *WorkflowTestSuite) TestWorkflow_Cancelled() {
	// Mock activity implementations.

	s.env.OnActivity(s.activities.Validate, mock.Anything, order.Order{}).Return(nil)

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow("cancelOrder", nil)
	}, time.Minute)

	// Execute workflow.

	s.env.ExecuteWorkflow(workflows.ProcessOrder, workflows.Params{Order: order.Order{}})

	// Assert execution and order status.

	s.Require().NoError(s.env.GetWorkflowError())

	val, err := s.env.QueryWorkflow("GetOrderStatus")
	s.Require().NoError(err, "workflow should be queryable")
	var got order.OrderStatus
	err = val.Get(&got)
	s.Require().NoError(err, "query result should be an order.OrderStatus")
	s.Equal(order.Cancelled, got, "order should be cancelled")
}
