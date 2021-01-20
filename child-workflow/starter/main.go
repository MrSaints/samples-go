package main

import (
	"context"
	"log"
	"time"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"

	child_workflow "github.com/temporalio/samples-go/child-workflow"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// This workflow ID can be user business logic identifier as well.
	workflowID := "parent-workflow_" + uuid.New()
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "child-workflow",
	}

	workflowRun, err := c.ExecuteWorkflow(context.Background(), workflowOptions, child_workflow.SampleParentWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow",
		"WorkflowID", workflowRun.GetID(), "RunID", workflowRun.GetRunID())

	time.Sleep(time.Second)

	log.Println("Cancelling")

	v, err := c.QueryWorkflow(context.Background(), workflowRun.GetID(), "", "child-workflow-id")
	if err != nil {
		log.Fatalln(err)
	}

	var childWorkflowID string
	err = v.Get(&childWorkflowID)
	if err != nil {
		log.Fatalln("Failed to get child workflow ID", err)
	}

	err = c.CancelWorkflow(context.Background(), childWorkflowID, "")
	if err != nil {
		log.Fatalln("Failed to cancel child workflow", err)
	}

	err = c.CancelWorkflow(context.Background(), workflowRun.GetID(), "")
	if err != nil {
		log.Fatalln("Failed to cancel parent workflow", err)
	}

	err = c.SignalWorkflow(
		context.Background(),
		childWorkflowID,
		"",
		"unblock",
		nil,
	)
	if err != nil {
		log.Fatalln("Failed to signal child workflow", err)
	}

	// Synchronously wait for the workflow completion. Behind the scenes the SDK performs a long poll operation.
	// If you need to wait for the workflow completion from another process use
	// Client.GetWorkflow API to get an instance of a WorkflowRun.
	var result string
	err = workflowRun.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Failure getting workflow result", err)
	}
	log.Println("Workflow result: %v", "result", result)
}
