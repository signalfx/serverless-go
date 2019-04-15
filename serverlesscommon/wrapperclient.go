package serverlesscommon

import (
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/signalfx/golib/sfxclient"
	log "github.com/sirupsen/logrus"
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
		log.Errorf("no value for environment variable %s", sfxAuthToken)
	}
	if os.Getenv(sfxIngestEndpoint) != "" {
		if ingestURL, err := url.Parse(os.Getenv(sfxIngestEndpoint)); err == nil {
			if ingestURL, err = ingestURL.Parse("v2/datapoint"); err == nil {
				handlerFuncWrapperClient.DatapointEndpoint = ingestURL.String()
			} else {
				log.Errorf("error parsing ingest url path v2/datapoint: %+v", err)
			}
		} else {
			log.Errorf("error parsing url value %s of environment variable %s. %+v", os.Getenv(sfxIngestEndpoint), sfxIngestEndpoint, err)
		}
	}
	if os.Getenv(sfxSendTimeoutSeconds) != "" {
		if timeout, err := time.ParseDuration(strings.TrimSpace(os.Getenv(sfxSendTimeoutSeconds)) + "s"); err == nil {
			handlerFuncWrapperClient.Client.Timeout = timeout
		} else {
			log.Errorf("error parsing timeout value %s of environment variable %s. %+v", os.Getenv(sfxSendTimeoutSeconds), sfxSendTimeoutSeconds, err)
		}
	}
}
