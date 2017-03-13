package main

import (
	"context"
	"encoding/json"
	"time"

	"code.uber.internal/devexp/minions-client-go.git/client/cadence"
	"code.uber.internal/devexp/minions-client-go.git/cmd/samples"
	"github.com/Sirupsen/logrus"
)

type (
	sampleActivity struct{}
	sampleWorkflow struct{}

	cronJobParam struct {
		SequenceID  int
		TotalJobNum int
	}
)

/**
* Workflow activities need to implement interface cadence.Activity.
 */
func (a sampleActivity) Execute(context context.Context, input []byte) ([]byte, error) {
	jobParam := &cronJobParam{}
	err := json.Unmarshal(input, jobParam)
	if err != nil {
		logrus.Panicf("Marshalling failed with error: %+v", err)
	}

	logrus.Infof("Done cron job %d of %d.", jobParam.SequenceID, jobParam.TotalJobNum)
	return nil, nil
}

func (a sampleActivity) ActivityType() cadence.ActivityType {
	return cadence.ActivityType{Name: sampleActivityName}
}

func (w sampleWorkflow) Execute(ctx cadence.Context, input []byte) (result []byte, err error) {
	totalJobNum := 3
	for i := 1; i <= totalJobNum; i++ {
		// Cadence is not ready to support true cron jobs yet. For now, we use this timer based job.
		// We will update this sample code when Cadence fully support the cron job.
		future := cadence.NewTimer(ctx, time.Second)
		cadence.NewSelector(ctx).AddFuture(
			future,
			func(v interface{}, err error) {
				input, err := json.Marshal(&cronJobParam{i, totalJobNum})
				if err != nil {
					logrus.Panicf("Marshalling failed with error: %+v", err)
					return
				}
				cadence.ExecuteActivity(ctx, samples.ActivityParameters(sampleActivityTaskList, sampleActivityName, input))
			}).Select(ctx)
	}

	return nil, nil
}
