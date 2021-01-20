package child_workflow

import (
	"go.temporal.io/sdk/workflow"
)

// This sample workflow demonstrates how to use invoke child workflow from parent workflow execution.  Each child
// workflow execution is starting a new run and parent execution is notified only after the completion of last run.

// SampleParentWorkflow workflow definition
func SampleParentWorkflow(ctx workflow.Context) (string, error) {
	logger := workflow.GetLogger(ctx)

	var childWorkflowID string
	err := workflow.SetQueryHandler(ctx, "child-workflow-id", func(input []byte) (string, error) {
		return childWorkflowID, nil
	})
	if err != nil {
		return "", err
	}

	cwo := workflow.ChildWorkflowOptions{}
	ctx = workflow.WithChildOptions(ctx, cwo)

	childWorkflowFuture := workflow.ExecuteChildWorkflow(ctx, SampleChildWorkflow, "World")

	var childWorkflowExecution workflow.Execution
	err = childWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &childWorkflowExecution)
	if err != nil {
		return "", err
	}

	childWorkflowID = childWorkflowExecution.ID

	var result string
	err = childWorkflowFuture.Get(ctx, &result)
	if err != nil {
		logger.Error("Parent execution received child execution failure.", "Error", err)
		return "", err
	}
	logger.Info("Parent execution completed.", "Result", result)
	return result, nil
}
