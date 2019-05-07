package gatewayserver

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"

	activatorutil "github.com/knative/serving/pkg/activator/util"
	"github.com/rancher/rio-autoscaler/pkg/controllers/gateway"
	"github.com/rancher/rio-autoscaler/pkg/logger"
	"github.com/rancher/rio-autoscaler/types"
	riov1controller "github.com/rancher/rio/pkg/generated/controllers/rio.cattle.io/v1"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	maxRetries             = 18 // the sum of all retries would add up to 1 minute
	minRetryInterval       = 100 * time.Millisecond
	exponentialBackoffBase = 1.3
	RioNameHeader          = "X-Rio-ServiceName"
	RioNamespaceHeader     = "X-Rio-Namespace"
	RioPortHeader          = "X-Rio-ServicePort"
	RequestCountHTTPHeader = "knative-activator-num-retries"
)

func NewHandler(rContext *types.Context) Handler {
	return Handler{
		services: rContext.Rio.Rio().V1().Service(),
	}
}

type Handler struct {
	services riov1controller.ServiceController
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get(RioNameHeader)
	namespace := r.Header.Get(RioNamespaceHeader)
	port := r.Header.Get(RioPortHeader)

	rioSvc, err := h.services.Get(namespace, name, metav1.GetOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if rioSvc.Status.ObservedScale != nil && *rioSvc.Status.ObservedScale == 0 {
		rioSvc.Status.ObservedScale = nil
		t := metav1.NewTime(time.Now())
		rioSvc.Status.ScaleFromZeroTimestamp = &t
		logrus.Infof("Activating service %s to scale 1", rioSvc.Name)
		if _, err := h.services.Update(rioSvc); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
	}

	timer := time.After(time.Minute)

	endpointCh, endpointNotReady := gateway.EndpointChanMap.Load(fmt.Sprintf("%s.%s", name, namespace))
	if !endpointNotReady {
		serveFQDN(name, namespace, port, w, r)
		return
	}
	select {
	case <-timer:
		http.Error(w, "timeout waiting for endpoint to be active", http.StatusGatewayTimeout)
		return
	case _, ok := <-endpointCh.(chan struct{}):
		if !ok {
			serveFQDN(name, namespace, port, w, r)
			return
		}
	}
}

func serveFQDN(name, namespace, port string, w http.ResponseWriter, r *http.Request) {
	targetUrl := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s.%s.svc.cluster.local:%s", name, namespace, port),
		Path:   r.URL.Path,
	}
	r.URL = targetUrl
	r.URL.Host = targetUrl.Host
	r.Host = targetUrl.Host

	// todo: check if 503 is actually coming from application or envoy
	shouldRetry := retryStatus(http.StatusServiceUnavailable)
	backoffSettings := wait.Backoff{
		Duration: minRetryInterval,
		Factor:   exponentialBackoffBase,
		Steps:    maxRetries,
	}

	rt := NewRetryRoundTripper(activatorutil.AutoTransport, logger.SugaredLogger, backoffSettings, shouldRetry)
	httpProxy := proxy.NewUpgradeAwareHandler(targetUrl, rt, true, false, er)
	httpProxy.ServeHTTP(w, r)
}

var er = &errorResponder{}

type errorResponder struct {
}

func (e *errorResponder) Error(w http.ResponseWriter, req *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	if _, err := w.Write([]byte(err.Error())); err != nil {
		logrus.Errorf("error writing response: %v", err)
	}
}

// the following code was taken from knative activator https://github.com/knative/serving/blob/eb01fb7078b18f16ee8cd97a46c6e57bf65123a8/pkg/activator/util/transports.go
type retryCond func(*http.Response) bool

func retryStatus(status int) retryCond {
	return func(resp *http.Response) bool {
		return resp.StatusCode == status
	}
}

type retryRoundTripper struct {
	logger          *zap.SugaredLogger
	transport       http.RoundTripper
	backoffSettings wait.Backoff
	retryConditions []retryCond
}

// RetryRoundTripper retries a request on error or retry condition, using the given `retry` strategy
func NewRetryRoundTripper(rt http.RoundTripper, l *zap.SugaredLogger, b wait.Backoff, conditions ...retryCond) http.RoundTripper {
	return &retryRoundTripper{
		logger:          l,
		transport:       rt,
		backoffSettings: b,
		retryConditions: conditions,
	}
}

func (rrt *retryRoundTripper) RoundTrip(r *http.Request) (resp *http.Response, err error) {
	// The request body cannot be read multiple times for retries.
	// The workaround is to clone the request body into a byte reader
	// so the body can be read multiple times.
	if r.Body != nil {
		rrt.logger.Debugf("Wrapping body in a rewinder.")
		r.Body = NewRewinder(r.Body)
	}

	attempts := 0
	wait.ExponentialBackoff(rrt.backoffSettings, func() (bool, error) {
		rrt.logger.Debugf("Retrying")

		attempts++
		r.Header.Add(RequestCountHTTPHeader, strconv.Itoa(attempts))
		resp, err = rrt.transport.RoundTrip(r)

		if err != nil {
			rrt.logger.Errorf("Error making a request: %s", err)
			return false, nil
		}

		for _, retryCond := range rrt.retryConditions {
			if retryCond(resp) {
				resp.Body.Close()
				return false, nil
			}
		}
		return true, nil
	})

	if err == nil {
		rrt.logger.Infof("Finished after %d attempt(s). Response code: %d", attempts, resp.StatusCode)

		if resp.Header == nil {
			resp.Header = make(http.Header)
		}

		resp.Header.Add(RequestCountHTTPHeader, strconv.Itoa(attempts))
	} else {
		rrt.logger.Errorf("Failed after %d attempts. Last error: %v", attempts, err)
	}

	return
}

type rewinder struct {
	sync.Mutex
	rc io.ReadCloser
	rs io.ReadSeeker
}

// rewinder wraps a single-use `ReadCloser` into a `ReadCloser` that can be read multiple times
func NewRewinder(rc io.ReadCloser) io.ReadCloser {
	return &rewinder{rc: rc}
}

func (r *rewinder) Read(b []byte) (int, error) {
	r.Lock()
	defer r.Unlock()
	// On the first `Read()`, the contents of `rc` is read into a buffer `rs`.
	// This buffer is used for all subsequent reads
	if r.rs == nil {
		buf, err := ioutil.ReadAll(r.rc)
		if err != nil {
			return 0, err
		}
		r.rc.Close()

		r.rs = bytes.NewReader(buf)
	}

	return r.rs.Read(b)
}

func (r *rewinder) Close() error {
	r.Lock()
	defer r.Unlock()
	// Rewind the buffer on `Close()` for the next call to `Read`
	r.rs.Seek(0, io.SeekStart)

	return nil
}
