# SignalFx Go Google Cloud Function Wrapper
SignalFx Golang Google Cloud Function Wrapper.

## Usage
The SignalFx Go Google Cloud Function Wrapper is a wrapper around an Google Cloud Function Go function handler, used to instrument execution of the function and send metrics to SignalFx.

### Installation
To install run the command:

`$ go get https://github.com/signalfx/google-cloud-function-go`

#### Configuring the ingest endpoint

By default, this function wrapper will send to the `us0` realm. If you are
not in this realm you will need to set the `SIGNALFX_INGEST_ENDPOINT` environment
variable to the correct realm ingest endpoint (https://ingest.{REALM}.signalfx.com/v2/datapoint).
To determine what realm you are in, check your profile page in the SignalFx
web application (click the avatar in the upper right and click My Profile).

### Environment Variable
Set the SIGNALFX_AUTH_TOKEN environment variable with the appropriate SignalFx authentication token. Change the default 
values of the other variables accordingly if desired.

`SIGNALFX_AUTH_TOKEN=<SignalFx authentication token>`

`SIGNALFX_INGEST_ENDPOINT=https://ingest.{REALM}.signalfx.com/v2/datapoint`

`SIGNALFX_SEND_TIMEOUT_SECONDS=5`

###  Wrapping a function
The SignalFx Go Google Cloud Function Wrapper wraps the handler `lambda.Handler`. Use the `lambda.NewHandler()` function to create the 
handler by passing your Lambda handler function to `lambda.NewHandler()`. Pass the created handler to the 
`sfxlambda.NewHandlerWrapper` function to create the wrapper `sfxlambda.HandlerWrapper`. Finally, pass the created wrapper 
to the `sfxlambda.Start()` function. See the example below.

```
import (
  ...
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/signalfx/lambda-go"
  ...
)
...

func handler(...) ... {
  ...  
}
...

func main() {
  ...
  handlerWrapper := sfxlambda.NewHandlerWrapper(lambda.NewHandler(handler))
  sfxlambda.Start(handlerWrapper)
  ...
}
...
```

### Metrics and dimensions sent by the wrapper
The Lambda wrapper sends the following metrics to SignalFx:

| Metric Name  | Type | Description |
| ------------- | ------------- | ---|
| function.invocations  | Counter  | Count number of Cloud Function invocations|
| function.cold_starts  | Counter  | Count number of cold starts|
| function.duration  | Gauge  | Milliseconds in execution time of underlying Cloud Function handler|

The Cloud Function wrapper adds the following dimensions to all data points sent to SignalFx:

| Dimension | Description |
| ------------- | ---|
| gcf_region  | Google Cloud Function Region  |
| gcf_project_id | Google Cloud Function Project ID  |
| gcf_function_name  | Google Cloud Function Name |
| gcf_function_version  | Google Cloud Function Version |=
| function_wrapper_version  | SignalFx function wrapper qualifier (e.g. signalfx_gcf_go-0.0.1) |
| metric_source | The literal value of 'gcf_wrapper' |


### Sending custom metric in the Google Cloud function
Use the method `sfxgcf.SendDatapoint()` of `HandlerWrapper` to send custom metric datapoints to SignalFx from within your 
Lambda handler function. A `sfxlambda.HandlerWrapper` variable needs to be declared globally in order to be accessible 
from within your Lambda handler function. See example below.

```
import (
  ...
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/signalfx/lambda-go"
  ...
)
...

var handlerWrapper sfxlambda.HandlerWrapper
...

func handler(w http.ResponseWriter, r *http.Request) {
  ...  
  // Custom counter metric.
  dp := datapoint.Datapoint {
      Metric: "db_calls",
      Value: datapoint.NewIntValue(1),
      MetricType: datapoint.Counter,
      Dimensions: map[string]string{"db_name":"mysql1",},
  }
  // Sending custom metric to SignalFx.
  handlerWrapper.SendDatapoints([]*datapoint.Datapoint{&dp})
  ...
}
...

func main(w http.ResponseWriter, r *http.Request) {
  ...
  handlerWrapper := NewHandlerWrapper(handler)
  handlerWrapper.Invoke(w, r)
  ...
}
...
```

### Testing locally.
WIP ~~Run the command below in the lambda-go package folder~~

WIP ~~`$ SIGNALFX_AUTH_TOKEN=test go test -v`~~

## License

Apache Software License v2. Copyright © 2014-2019 SignalFx
