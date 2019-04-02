package sfxserverlesscommon

import (
	"context"
	"fmt"
	"github.com/signalfx/golib/datapoint"
	"time"
)

var SendDatapoints = func(ctx context.Context, dps []*datapoint.Datapoint) error {
	now := time.Now()
	for _, dp := range dps {
		if dp.Timestamp.IsZero() {
			dp.Timestamp = now
		}
	}
	if err := handlerFuncWrapperClient.AddDatapoints(ctx, dps); err != nil {
		return fmt.Errorf("error sending datapoint to SignalFx. %+v", err)
	}
	return nil
}

func InvocationsDatapoint() *datapoint.Datapoint {
	dp := datapoint.Datapoint{Metric: "function.invocations", Value: datapoint.NewIntValue(1), MetricType: datapoint.Counter}
	return &dp
}

func ColdStartsDatapoint() *datapoint.Datapoint {
	dp := datapoint.Datapoint{Metric: "function.cold_starts", Value: datapoint.NewIntValue(1), MetricType: datapoint.Counter}
	return &dp
}

func DurationDatapoint(elapsed time.Duration) *datapoint.Datapoint {
	dp := datapoint.Datapoint{Metric: "function.duration", Value: datapoint.NewFloatValue(elapsed.Seconds()), MetricType: datapoint.Gauge}
	return &dp
}

func ErrorsDatapoint() *datapoint.Datapoint {
	dp := datapoint.Datapoint{Metric: "function.errors", Value: datapoint.NewIntValue(1), MetricType: datapoint.Counter}
	return &dp
}
