package sfxgcf

import (
	"context"
	"fmt"
	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/golib/sfxclient"
	"net/url"
	"os"
	"strings"
	"time"
)

var handlerFuncWrapperClient *sfxclient.HTTPSink

const (
	sfxAuthToken          = "SIGNALFX_AUTH_TOKEN"
	sfxIngestEndpoint     = "SIGNALFX_INGEST_ENDPOINT"
	sfxSendTimeoutSeconds = "SIGNALFX_SEND_TIMEOUT_SECONDS"
)

func init() {
	handlerFuncWrapperClient = sfxclient.NewHTTPSink()
	if handlerFuncWrapperClient.AuthToken = os.Getenv(sfxAuthToken); handlerFuncWrapperClient.AuthToken == "" {
		logger.error.Printf("no value for environment variable %s", sfxAuthToken)
	}
	if os.Getenv(sfxIngestEndpoint) != "" {
		if ingestURL, err := url.Parse(os.Getenv(sfxIngestEndpoint)); err == nil {
			if ingestURL, err = ingestURL.Parse("v2/datapoint"); err == nil {
				handlerFuncWrapperClient.DatapointEndpoint = ingestURL.String()
			} else {
				logger.error.Printf("error parsing ingest url path v2/datapoint: %+v", err)
			}
		} else {
			logger.error.Printf("error parsing url value %s of environment variable %s. %+v", os.Getenv(sfxIngestEndpoint), sfxIngestEndpoint, err)
		}
	}
	if os.Getenv(sfxSendTimeoutSeconds) != "" {
		if timeout, err := time.ParseDuration(strings.TrimSpace(os.Getenv(sfxSendTimeoutSeconds)) + "s"); err == nil {
			handlerFuncWrapperClient.Client.Timeout = timeout
		} else {
			logger.error.Printf("error parsing timeout value %s of environment variable %s. %+v", os.Getenv(sfxSendTimeoutSeconds), sfxSendTimeoutSeconds, err)
		}
	}
}

var sendDatapoints = func(ctx context.Context, dps []*datapoint.Datapoint) error {
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

func invocationsDatapoint() *datapoint.Datapoint {
	dp := datapoint.Datapoint{Metric: "function.invocations", Value: datapoint.NewIntValue(1), MetricType: datapoint.Counter}
	return &dp
}

func coldStartsDatapoint() *datapoint.Datapoint {
	dp := datapoint.Datapoint{Metric: "function.cold_starts", Value: datapoint.NewIntValue(1), MetricType: datapoint.Counter}
	return &dp
}

func durationDatapoint(elapsed time.Duration) *datapoint.Datapoint {
	dp := datapoint.Datapoint{Metric: "function.duration", Value: datapoint.NewFloatValue(elapsed.Seconds()), MetricType: datapoint.Gauge}
	return &dp
}

func errorsDatapoint() *datapoint.Datapoint {
	dp := datapoint.Datapoint{Metric: "function.errors", Value: datapoint.NewIntValue(1), MetricType: datapoint.Counter}
	return &dp
}