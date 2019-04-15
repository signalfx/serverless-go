package serverlesscommon

import (
	"context"
	"fmt"
	"time"

	"github.com/signalfx/golib/datapoint"
)

func SendDatapoints(ctx context.Context, dps []*datapoint.Datapoint) error {
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

// InvocationsDatapoint creates a datapoint to report function invocations count to SignalFx
func InvocationsDatapoint() *datapoint.Datapoint {
	dp := datapoint.Datapoint{Metric: "function.invocations", Value: datapoint.NewIntValue(1), MetricType: datapoint.Counter}
	return &dp
}

// ColdStartsDatapoint creates a datapoint to report function cold starts count to SignalFx
func ColdStartsDatapoint() *datapoint.Datapoint {
	dp := datapoint.Datapoint{Metric: "function.cold_starts", Value: datapoint.NewIntValue(1), MetricType: datapoint.Counter}
	return &dp
}

// DurationDatapoint creates a datapoint to report function duration to SignalFx
func DurationDatapoint(elapsed time.Duration) *datapoint.Datapoint {
	dp := datapoint.Datapoint{Metric: "function.duration", Value: datapoint.NewFloatValue(elapsed.Seconds()), MetricType: datapoint.Gauge}
	return &dp
}

// ErrorsDatapoint creates a datapoint to report function errors count to SignalFx
func ErrorsDatapoint() *datapoint.Datapoint {
	dp := datapoint.Datapoint{Metric: "function.errors", Value: datapoint.NewIntValue(1), MetricType: datapoint.Counter}
	return &dp
}
