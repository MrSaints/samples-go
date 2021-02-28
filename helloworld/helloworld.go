package helloworld

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

// Workflow is a Hello World workflow definition.
func Workflow(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", name)

	var uniqueID *string

	v := workflow.GetVersion(ctx, "UniqueID", workflow.DefaultVersion, 1)
	if v == 1 {
		encodedUID := workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
			return "TEST-UNIQUE-ID"
		})
		err := encodedUID.Get(&uniqueID)
		if err != nil {
			logger.Error("Side effect failed.", "Error", err)
			return "", err
		}
	}

	var result string
	err := workflow.ExecuteActivity(ctx, Activity, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("HelloWorld workflow completed.", "result", result)

	return result, nil
}

func Activity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", name)
	return "Hello " + name + "!", nil
}
