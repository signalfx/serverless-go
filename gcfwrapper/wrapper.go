package gcfwrapper

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/signalfx/golib/datapoint"

	sfxcommon "github.com/seonsfx/serverless-go/serverlesscommon"
	log "github.com/sirupsen/logrus"
)

const (
	name    = "signalfx_gcf_go"
	version = "0.0.1"
)

// HandlerWrapper provides methods to invoke user function and send custom datapoints
type HandlerWrapper struct {
	handler      func(http.ResponseWriter, *http.Request)
	ctx          context.Context
	notColdStart bool
}

// NewHandlerWrapper is a HandlerWrapper creating factory function.
func NewHandlerWrapper(h func(http.ResponseWriter, *http.Request)) HandlerWrapper {
	return HandlerWrapper{handler: h}
}

// Invoke runs user function and records the invocation/duration of function
func (hw *HandlerWrapper) Invoke(w http.ResponseWriter, r *http.Request) {
	hw.ctx = r.Context()
	dps := []*datapoint.Datapoint{sfxcommon.InvocationsDatapoint()}
	if !hw.notColdStart {
		dps = append(dps, sfxcommon.ColdStartsDatapoint())
		hw.notColdStart = true
	}
	start := time.Now()
	hw.handler(w, r)
	dps = append(dps, sfxcommon.DurationDatapoint(time.Since(start)))

	if err := hw.SendDatapoints(dps); err != nil {
		log.Error(err)
	}
}

// SendDatapoints sends custom metric datapoints to SignalFx.
func (hw *HandlerWrapper) SendDatapoints(dps []*datapoint.Datapoint) error {
	if hw.ctx == nil {
		return fmt.Errorf("invalid argument. request is nil")
	}

	dims := defaultDimensions(hw.ctx)

	for _, dp := range dps {
		dp.Dimensions = datapoint.AddMaps(dims, dp.Dimensions)
	}

	if err := sfxcommon.SendDatapoints(hw.ctx, dps); err != nil {
		return err
	}

	return nil
}

type dimensions map[string]string

func defaultDimensions(ctx context.Context) map[string]string {
	dims := dimensions{
		"metric_source":            "gcf_wrapper",
		"function_wrapper_version": name + "_" + version,
	}
	if os.Getenv("FUNCTION_REGION") != "" {
		dims["gcf_region"] = os.Getenv("FUNCTION_REGION")
	}
	if os.Getenv("GCP_PROJECT") != "" {
		dims["gcf_project_id"] = os.Getenv("GCP_PROJECT")
	}
	if os.Getenv("FUNCTION_NAME") != "" {
		dims["gcf_function_name"] = os.Getenv("FUNCTION_NAME")
	}
	if os.Getenv("X_GOOGLE_FUNCTION_VERSION") != "" {
		dims["gcf_function_version"] = os.Getenv("X_GOOGLE_FUNCTION_VERSION")
	}

	return dims
}

// InvocationsDatapoint exposes a function from common library to create a datapoint to report function invocations count to SignalFx
var InvocationsDatapoint = sfxcommon.InvocationsDatapoint

// ColdStartsDatapoint exposes a function from common library to create a datapoint to report function cold starts count to SignalFx
var ColdStartsDatapoint = sfxcommon.ColdStartsDatapoint

// DurationDatapoint exposes a function from common library to create a datapoint to report function duration to SignalFx
var DurationDatapoint = sfxcommon.DurationDatapoint

// ErrorsDatapoint exposes a function from common library to create a datapoint to report function errors count to SignalFx
var ErrorsDatapoint = sfxcommon.ErrorsDatapoint
