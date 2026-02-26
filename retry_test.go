package slack

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// TestRetryOn429ThenSuccess verifies that 429 responses are retried (using Retry-After or config)
// until success; call count = 2×429 + 1×200 = 3.
func TestRetryOn429ThenSuccess(t *testing.T) {
	t.Parallel()

	var callCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n <= 2 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 3
	cfg.RetryAfterDuration = 1 * time.Millisecond
	cfg.BackoffInitial = time.Millisecond
	cfg.BackoffMax = 5 * time.Millisecond
	cfg.Handlers = AllBuiltinRetryHandlers(cfg)

	api := New("token", OptionAPIURL(srv.URL+"/"), OptionRetryConfig(cfg))
	var out SlackResponse
	err := api.postMethod(context.Background(), "auth.test", nil, &out)
	if err != nil {
		t.Fatalf("postMethod: %v", err)
	}
	if !out.Ok {
		t.Errorf("want ok=true, got ok=%v", out.Ok)
	}
	if got := atomic.LoadInt32(&callCount); got != 3 {
		t.Errorf("want 3 calls (2x 429 + 1x 200), got %d", got)
	}
}

// TestRetryOn500ThenSuccess verifies that 5xx responses are retried with backoff until success.
func TestRetryOn500ThenSuccess(t *testing.T) {
	t.Parallel()

	var callCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 3
	cfg.BackoffInitial = time.Millisecond
	cfg.BackoffMax = 5 * time.Millisecond
	cfg.Handlers = append(AllBuiltinRetryHandlers(cfg), NewServerErrorRetryHandler(cfg))

	api := New("token", OptionAPIURL(srv.URL+"/"), OptionRetryConfig(cfg))
	var out SlackResponse
	err := api.postMethod(context.Background(), "auth.test", nil, &out)
	if err != nil {
		t.Fatalf("postMethod: %v", err)
	}
	if !out.Ok {
		t.Errorf("want ok=true, got ok=%v", out.Ok)
	}
	if got := atomic.LoadInt32(&callCount); got != 3 {
		t.Errorf("want 3 calls (2x 500 + 1x 200), got %d", got)
	}
}

// TestRetryExhaustedReturnsLastError verifies that when all 5xx retries are exhausted we return
// StatusCodeError and the server was called MaxRetries+1 times (3 for MaxRetries=2).
func TestRetryExhaustedReturnsLastError(t *testing.T) {
	t.Parallel()

	var callCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 2
	cfg.BackoffInitial = time.Millisecond
	cfg.BackoffMax = time.Millisecond
	cfg.Handlers = append(AllBuiltinRetryHandlers(cfg), NewServerErrorRetryHandler(cfg))

	api := New("token", OptionAPIURL(srv.URL+"/"), OptionRetryConfig(cfg))
	var out SlackResponse
	err := api.postMethod(context.Background(), "auth.test", nil, &out)
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	if _, ok := err.(StatusCodeError); !ok {
		t.Errorf("expected StatusCodeError, got %T: %v", err, err)
	}
	if got := atomic.LoadInt32(&callCount); got != 3 {
		t.Errorf("want 3 calls (all 500), got %d", got)
	}
}

// TestRetryExhausted5xxResponseBodyReadable verifies that when retries are exhausted on 5xx,
// the response body is left readable so checkStatusCode/logResponse can dump it for debug.
func TestRetryExhausted5xxResponseBodyReadable(t *testing.T) {
	t.Parallel()

	const bodyContent = "custom error body from server"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(bodyContent))
	}))
	defer srv.Close()

	buf := bytes.NewBuffer(nil)
	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 2
	cfg.BackoffInitial = time.Millisecond
	cfg.BackoffMax = time.Millisecond
	cfg.Handlers = append(AllBuiltinRetryHandlers(cfg), NewServerErrorRetryHandler(cfg))

	api := New("token",
		OptionAPIURL(srv.URL+"/"),
		OptionRetryConfig(cfg),
		OptionDebug(true),
		OptionLog(log.New(buf, "", 0)),
	)
	var out SlackResponse
	err := api.postMethod(context.Background(), "auth.test", nil, &out)
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	if _, ok := err.(StatusCodeError); !ok {
		t.Errorf("expected StatusCodeError, got %T: %v", err, err)
	}
	logged := buf.String()
	if !strings.Contains(logged, bodyContent) {
		t.Errorf("debug log should contain response body %q (so checkStatusCode/logResponse could read it); got: %s", bodyContent, logged)
	}
}

// TestRetryExhausted429ResponseBodyReadable verifies that when retries are exhausted on 429
// without Retry-After, the response body is left readable so logResponse can dump it for debug.
func TestRetryExhausted429ResponseBodyReadable(t *testing.T) {
	t.Parallel()

	const bodyContent = "rate limit message from server"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(bodyContent))
	}))
	defer srv.Close()

	buf := bytes.NewBuffer(nil)
	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 2
	cfg.RetryAfterDuration = time.Millisecond
	cfg.RetryAfterJitter = 0
	cfg.Handlers = AllBuiltinRetryHandlers(cfg)

	api := New("token",
		OptionAPIURL(srv.URL+"/"),
		OptionRetryConfig(cfg),
		OptionDebug(true),
		OptionLog(log.New(buf, "", 0)),
	)
	var out SlackResponse
	err := api.postMethod(context.Background(), "auth.test", nil, &out)
	if err == nil {
		t.Fatal("expected error after exhausting 429 retries")
	}
	// 429 without Retry-After returns StatusCodeError and calls logResponse
	if _, ok := err.(StatusCodeError); !ok {
		t.Errorf("expected StatusCodeError (429 without Retry-After), got %T: %v", err, err)
	}
	logged := buf.String()
	if !strings.Contains(logged, bodyContent) {
		t.Errorf("debug log should contain response body %q; got: %s", bodyContent, logged)
	}
}

// TestRetryExhausted429ReturnsError verifies that when 429 retries are exhausted (with Retry-After
// set) we return *RateLimitedError and the server was called MaxRetries+1 times.
func TestRetryExhausted429ReturnsError(t *testing.T) {
	t.Parallel()

	var callCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 2
	cfg.RetryAfterDuration = time.Millisecond
	cfg.RetryAfterJitter = 0
	cfg.Handlers = AllBuiltinRetryHandlers(cfg)

	api := New("token", OptionAPIURL(srv.URL+"/"), OptionRetryConfig(cfg))
	var out SlackResponse
	err := api.postMethod(context.Background(), "auth.test", nil, &out)
	if err == nil {
		t.Fatal("expected error after exhausting 429 retries")
	}
	if _, ok := err.(*RateLimitedError); !ok {
		t.Errorf("expected *RateLimitedError, got %T: %v", err, err)
	}
	if got := atomic.LoadInt32(&callCount); got != 3 {
		t.Errorf("want 3 calls (all 429), got %d", got)
	}
}

// TestOptionRetryZeroDisablesRetry verifies OptionRetryConfig(RetryConfig{MaxRetries: 0}) does not
// wrap the client (no retry layer).
func TestOptionRetryZeroDisablesRetry(t *testing.T) {
	t.Parallel()

	api := New("token", OptionRetryConfig(RetryConfig{MaxRetries: 0}))
	if _, ok := api.httpclient.(*retryClient); ok {
		t.Error("OptionRetryConfig with MaxRetries=0 should not wrap client")
	}
}

// TestOptionRetryNonPositiveDisablesRetry verifies OptionRetry(0), OptionRetry(-1), etc. do not
// wrap the client.
func TestOptionRetryNonPositiveDisablesRetry(t *testing.T) {
	t.Parallel()

	for _, maxRetries := range []int{0, -1, -10} {
		api := New("token", OptionRetry(maxRetries))
		if _, ok := api.httpclient.(*retryClient); ok {
			t.Errorf("OptionRetry(%d) should not wrap client", maxRetries)
		}
	}
}

// TestRetryOn503ThenSuccess verifies 503 is retried like other 5xx until success.
func TestRetryOn503ThenSuccess(t *testing.T) {
	t.Parallel()

	var callCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 3
	cfg.BackoffInitial = time.Millisecond
	cfg.BackoffMax = 5 * time.Millisecond
	cfg.Handlers = append(AllBuiltinRetryHandlers(cfg), NewServerErrorRetryHandler(cfg))

	api := New("token", OptionAPIURL(srv.URL+"/"), OptionRetryConfig(cfg))
	var out SlackResponse
	err := api.postMethod(context.Background(), "auth.test", nil, &out)
	if err != nil {
		t.Fatalf("postMethod: %v", err)
	}
	if !out.Ok {
		t.Errorf("want ok=true, got ok=%v", out.Ok)
	}
	if got := atomic.LoadInt32(&callCount); got != 3 {
		t.Errorf("want 3 calls (2x 503 + 1x 200), got %d", got)
	}
}

// TestRetryOnConnectionErrorThenSuccess verifies connection errors (e.g. "connection reset") are
// retried until success; call count = 2×error + 1×200 = 3.
func TestRetryOnConnectionErrorThenSuccess(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	base := &http.Client{}
	wrapped := &failingThenOKClient{httpClient: base, failAttempts: 2}

	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 3
	cfg.BackoffInitial = time.Millisecond
	cfg.BackoffMax = 5 * time.Millisecond
	cfg.Handlers = AllBuiltinRetryHandlers(cfg) // connection + 429; default is 429 only

	api := New("token", OptionAPIURL(srv.URL+"/"), OptionHTTPClient(wrapped), OptionRetryConfig(cfg))
	var out SlackResponse
	err := api.postMethod(context.Background(), "auth.test", nil, &out)
	if err != nil {
		t.Fatalf("postMethod: %v", err)
	}
	if !out.Ok {
		t.Errorf("want ok=true, got ok=%v", out.Ok)
	}
	if got := atomic.LoadInt32(&wrapped.attempts); got != 3 {
		t.Errorf("want 3 calls (2x connection error + 1x 200), got %d", got)
	}
}

// TestRetryOnConnectionErrorExhausted verifies that when connection errors persist for all attempts
// we return the last error (no response) and the underlying client was called MaxRetries+1 times.
func TestRetryOnConnectionErrorExhausted(t *testing.T) {
	t.Parallel()

	failClient := &connectionFailingClient{}
	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 2
	cfg.BackoffInitial = time.Millisecond
	cfg.BackoffMax = time.Millisecond
	cfg.Handlers = AllBuiltinRetryHandlers(cfg) // connection + 429; default is 429 only

	rc := &retryClient{client: failClient, config: cfg, debug: nil}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://localhost/", strings.NewReader("body"))
	resp, err := rc.Do(req)
	if err == nil {
		t.Fatal("expected error when all connection attempts fail")
	}
	if resp != nil {
		t.Errorf("expected nil response when error returned, got %v", resp)
	}
	if got := atomic.LoadInt32(&failClient.attempts); got != 3 {
		t.Errorf("want 3 attempts (MaxRetries+1), got %d", got)
	}
}

// TestNoRetryWhenGetBodyNil verifies that when a request has a body but no GetBody we do not
// retry on 429/5xx, so we never send an empty body; the underlying client is called exactly once.
func TestNoRetryWhenGetBodyNil(t *testing.T) {
	t.Parallel()

	mock := &countAndRespondClient{code: http.StatusTooManyRequests}
	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 3
	cfg.RetryAfterDuration = time.Millisecond
	cfg.RetryAfterJitter = 0

	rc := &retryClient{client: mock, config: cfg, debug: nil}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://localhost/", io.NopCloser(strings.NewReader("body")))
	if err != nil {
		t.Fatal(err)
	}
	req.GetBody = nil // request body cannot be replayed; retry must not send empty body

	resp, err := rc.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	resp.Body.Close()
	if got := atomic.LoadInt32(&mock.count); got != 1 {
		t.Errorf("when body present and GetBody is nil should not retry, got %d calls", got)
	}
}

// TestRetryGetRequestWithNilBody verifies that GET requests (or any request with nil body and
// no GetBody, as created by getResource) are retried on 429 when using AllBuiltinRetryHandlers.
func TestRetryGetRequestWithNilBody(t *testing.T) {
	t.Parallel()

	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 3
	cfg.RetryAfterDuration = time.Millisecond
	cfg.RetryAfterJitter = 0
	cfg.Handlers = AllBuiltinRetryHandlers(cfg)

	// Return 429 twice then 200 on third call.
	wrapped := &retryCountThenSuccessClient{needFail: 2}
	rc := &retryClient{client: wrapped, config: cfg, debug: nil}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Body and GetBody are both nil, as with getResource — should still retry.
	resp, err := rc.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	resp.Body.Close()
	if got := atomic.LoadInt32(&wrapped.calls); got != 3 {
		t.Errorf("GET with nil body should retry on 429 (2x 429 + 1x 200), got %d calls", got)
	}
}

// contextCancelClient returns context.Canceled when the request context is already cancelled.
type contextCancelClient struct {
	httpClient
	calls int32
}

func (c *contextCancelClient) Do(req *http.Request) (*http.Response, error) {
	atomic.AddInt32(&c.calls, 1)
	if err := req.Context().Err(); err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}

// TestRetryRespectsContextCancelation verifies that when the request context is already
// cancelled before the first attempt, we return the context error without retrying
// (underlying client is called once and gets context.Canceled).
func TestRetryRespectsContextCancelation(t *testing.T) {
	t.Parallel()

	mock := &countAndRespondClient{code: http.StatusTooManyRequests}
	wrapped := &contextCancelClient{httpClient: mock}
	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 2
	cfg.RetryAfterDuration = 60 * time.Second
	cfg.RetryAfterJitter = 0
	cfg.Handlers = AllBuiltinRetryHandlers(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before first attempt

	rc := &retryClient{client: wrapped, config: cfg, debug: nil}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost/", strings.NewReader("body"))
	if err != nil {
		t.Fatal(err)
	}
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader("body")), nil
	}
	resp, err := rc.Do(req)
	if err == nil {
		t.Fatal("expected error when context is cancelled")
	}
	if resp != nil {
		resp.Body.Close()
		t.Errorf("expected nil response when context cancelled, got %v", resp)
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
	if got := atomic.LoadInt32(&wrapped.calls); got != 1 {
		t.Errorf("expected 1 call, got %d", got)
	}
}

// TestRetryContextCanceledDuringSleep verifies that when the context is cancelled during
// the retry wait (after a 429), we return immediately with the context error instead of
// sleeping the full duration, and we do not perform a second request.
func TestRetryContextCanceledDuringSleep(t *testing.T) {
	t.Parallel()

	mock := &countAndRespondClient{code: http.StatusTooManyRequests}
	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 2
	cfg.RetryAfterDuration = 60 * time.Second // long sleep; cancel during this
	cfg.RetryAfterJitter = 0
	cfg.Handlers = AllBuiltinRetryHandlers(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	body := "body"
	rc := &retryClient{client: mock, config: cfg, debug: nil}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader(body)), nil
	}

	// Cancel context shortly after first attempt returns 429, so we exit during sleep.
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	resp, err := rc.Do(req)
	if err == nil {
		t.Fatal("expected error when context is cancelled during sleep")
	}
	if resp != nil {
		resp.Body.Close()
		t.Errorf("expected nil response when context cancelled, got %v", resp)
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
	if got := atomic.LoadInt32(&mock.count); got != 1 {
		t.Errorf("expected 1 call (cancel during sleep before retry), got %d", got)
	}
}

// TestDefaultRetryConfig verifies DefaultRetryConfig returns the documented defaults.
func TestDefaultRetryConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultRetryConfig()
	if cfg.MaxRetries != 3 {
		t.Errorf("MaxRetries: got %d, want 3", cfg.MaxRetries)
	}
	if cfg.RetryAfterDuration != 60*time.Second {
		t.Errorf("RetryAfterDuration: got %v, want 60s", cfg.RetryAfterDuration)
	}
	if cfg.RetryAfterJitter != 1*time.Second {
		t.Errorf("RetryAfterJitter: got %v, want 1s", cfg.RetryAfterJitter)
	}
	if cfg.BackoffInitial != 100*time.Millisecond {
		t.Errorf("BackoffInitial: got %v, want 100ms", cfg.BackoffInitial)
	}
	if cfg.BackoffMax != 30*time.Second {
		t.Errorf("BackoffMax: got %v, want 30s", cfg.BackoffMax)
	}
	if cfg.BackoffJitter != 50*time.Millisecond {
		t.Errorf("BackoffJitter: got %v, want 50ms", cfg.BackoffJitter)
	}
}

// TestConnectionOnlyRetryHandlersOnlyRetriesConnection verifies that ConnectionOnlyRetryHandlers
// (connection only) does not retry 429 or 5xx; one call each, then error.
func TestConnectionOnlyRetryHandlersOnlyRetriesConnection(t *testing.T) {
	t.Parallel()

	t.Run("429_not_retried", func(t *testing.T) {
		t.Parallel()
		var callCount int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&callCount, 1)
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer srv.Close()
		cfg := DefaultRetryConfig()
		cfg.MaxRetries = 3
		cfg.Handlers = ConnectionOnlyRetryHandlers()
		api := New("token", OptionAPIURL(srv.URL+"/"), OptionRetryConfig(cfg))
		var out SlackResponse
		err := api.postMethod(context.Background(), "auth.test", nil, &out)
		if err == nil {
			t.Fatal("expected error on 429 with connection-only handlers")
		}
		if got := atomic.LoadInt32(&callCount); got != 1 {
			t.Errorf("429 should not be retried with ConnectionOnlyRetryHandlers, got %d calls", got)
		}
	})

	t.Run("5xx_not_retried", func(t *testing.T) {
		t.Parallel()
		var callCount int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&callCount, 1)
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()
		cfg := DefaultRetryConfig()
		cfg.MaxRetries = 3
		cfg.Handlers = ConnectionOnlyRetryHandlers()
		api := New("token", OptionAPIURL(srv.URL+"/"), OptionRetryConfig(cfg))
		var out SlackResponse
		err := api.postMethod(context.Background(), "auth.test", nil, &out)
		if err == nil {
			t.Fatal("expected error on 500 with connection-only handlers")
		}
		if got := atomic.LoadInt32(&callCount); got != 1 {
			t.Errorf("5xx should not be retried with ConnectionOnlyRetryHandlers, got %d calls", got)
		}
	})
}

// TestOptionRetryRetries429Not5xx verifies OptionRetry(n) uses DefaultRetryHandlers (429 only): 429 retried, 5xx not retried.
func TestOptionRetryRetries429Not5xx(t *testing.T) {
	t.Parallel()

	t.Run("429_retried", func(t *testing.T) {
		t.Parallel()
		var callCount int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			n := atomic.AddInt32(&callCount, 1)
			if n <= 2 {
				w.Header().Set("Retry-After", "0")
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"ok":true}`))
		}))
		defer srv.Close()
		api := New("token", OptionAPIURL(srv.URL+"/"), OptionRetry(3))
		var out SlackResponse
		err := api.postMethod(context.Background(), "auth.test", nil, &out)
		if err != nil {
			t.Fatalf("postMethod: %v", err)
		}
		if got := atomic.LoadInt32(&callCount); got != 3 {
			t.Errorf("OptionRetry should retry 429, want 3 calls, got %d", got)
		}
	})

	t.Run("5xx_not_retried", func(t *testing.T) {
		t.Parallel()
		var callCount int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&callCount, 1)
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()
		api := New("token", OptionAPIURL(srv.URL+"/"), OptionRetry(3))
		var out SlackResponse
		err := api.postMethod(context.Background(), "auth.test", nil, &out)
		if err == nil {
			t.Fatal("expected error on 500 with OptionRetry (no 5xx handler)")
		}
		if got := atomic.LoadInt32(&callCount); got != 1 {
			t.Errorf("OptionRetry should not retry 5xx by default, got %d calls", got)
		}
	})
}

// TestOptionRetryConfigWithNilHandlersDefaultsTo429Only verifies that OptionRetryConfig(cfg)
// with cfg.Handlers == nil defaults to DefaultRetryHandlers (429 only); 429 is retried until success.
func TestOptionRetryConfigWithNilHandlersDefaultsTo429Only(t *testing.T) {
	t.Parallel()

	var callCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n <= 2 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()
	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 3
	// Handlers explicitly nil; OptionRetryConfig sets DefaultRetryHandlers(cfg) (429 only).
	api := New("token", OptionAPIURL(srv.URL+"/"), OptionRetryConfig(cfg))
	var out SlackResponse
	err := api.postMethod(context.Background(), "auth.test", nil, &out)
	if err != nil {
		t.Fatalf("postMethod: %v", err)
	}
	if !out.Ok {
		t.Errorf("want ok=true, got ok=%v", out.Ok)
	}
	if got := atomic.LoadInt32(&callCount); got != 3 {
		t.Errorf("with nil Handlers default is DefaultRetryHandlers (429 only), 429 should retry; want 3 calls, got %d", got)
	}
}

// TestOptionRetryConfigWithNilHandlersDoesNotRetryConnection verifies that with default (429 only)
// handlers, connection errors are not retried; one call then error.
func TestOptionRetryConfigWithNilHandlersDoesNotRetryConnection(t *testing.T) {
	t.Parallel()

	base := &http.Client{}
	wrapped := &failingThenOKClient{httpClient: base, failAttempts: 2}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 3
	// Handlers nil = DefaultRetryHandlers (429 only); connection errors not retried
	api := New("token", OptionAPIURL(srv.URL+"/"), OptionHTTPClient(wrapped), OptionRetryConfig(cfg))
	var out SlackResponse
	err := api.postMethod(context.Background(), "auth.test", nil, &out)
	if err == nil {
		t.Fatal("expected error when connection fails with default (429-only) handlers")
	}
	if got := atomic.LoadInt32(&wrapped.attempts); got != 1 {
		t.Errorf("with default (429 only) handlers connection should not be retried, got %d calls", got)
	}
}

// TestServerErrorRetryHandlerOptIn verifies that adding NewServerErrorRetryHandler retries 5xx.
func TestServerErrorRetryHandlerOptIn(t *testing.T) {
	t.Parallel()

	var callCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()
	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 3
	cfg.BackoffInitial = time.Millisecond
	cfg.BackoffMax = 5 * time.Millisecond
	cfg.Handlers = append(AllBuiltinRetryHandlers(cfg), NewServerErrorRetryHandler(cfg))
	api := New("token", OptionAPIURL(srv.URL+"/"), OptionRetryConfig(cfg))
	var out SlackResponse
	err := api.postMethod(context.Background(), "auth.test", nil, &out)
	if err != nil {
		t.Fatalf("postMethod: %v", err)
	}
	if !out.Ok {
		t.Errorf("want ok=true, got ok=%v", out.Ok)
	}
	if got := atomic.LoadInt32(&callCount); got != 3 {
		t.Errorf("with ServerErrorRetryHandler want 3 calls (2x 500 + 1x 200), got %d", got)
	}
}

// TestIsRetryableConnError verifies which errors are treated as retryable connection failures.
func TestIsRetryableConnError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"connection reset", errConnectionReset, true},
		{"connection refused", errConnectionRefused, true},
		{"EOF", errEOF, true},
		{"other", errOther, false},
		{"wrapped connection reset", fmt.Errorf("request failed: %w", errConnectionReset), true},
		{"wrapped other", fmt.Errorf("request failed: %w", errOther), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := isRetryableConnError(tt.err)
			if got != tt.want {
				t.Errorf("isRetryableConnError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

// --- Test helpers ---

// failingThenOKClient fails with a retryable connection error the first N times, then delegates.
type failingThenOKClient struct {
	httpClient
	failAttempts int32
	attempts     int32
}

func (c *failingThenOKClient) Do(req *http.Request) (*http.Response, error) {
	n := atomic.AddInt32(&c.attempts, 1)
	if n <= c.failAttempts {
		return nil, &mockErr{msg: "connection reset by peer"}
	}
	return c.httpClient.Do(req)
}

// countAndRespondClient counts Do calls and returns a fixed HTTP status code.
type countAndRespondClient struct {
	count int32
	code  int
}

func (c *countAndRespondClient) Do(req *http.Request) (*http.Response, error) {
	atomic.AddInt32(&c.count, 1)
	return &http.Response{
		StatusCode: c.code,
		Body:       io.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
	}, nil
}

// retryCountThenSuccessClient returns 429 for the first needFail Do calls, then 200.
type retryCountThenSuccessClient struct {
	needFail int32
	calls    int32
}

func (c *retryCountThenSuccessClient) Do(req *http.Request) (*http.Response, error) {
	n := atomic.AddInt32(&c.calls, 1)
	code := http.StatusTooManyRequests
	if n > c.needFail {
		code = http.StatusOK
	}
	return &http.Response{
		StatusCode: code,
		Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
		Header:     make(http.Header),
	}, nil
}

// connectionFailingClient always returns a retryable connection error; attempts is incremented each Do.
type connectionFailingClient struct {
	attempts int32
}

func (c *connectionFailingClient) Do(*http.Request) (*http.Response, error) {
	atomic.AddInt32(&c.attempts, 1)
	return nil, &mockErr{msg: "connection reset by peer"}
}

var (
	errConnectionReset   = &mockErr{msg: "connection reset by peer"}
	errConnectionRefused = &mockErr{msg: "connection refused"}
	errEOF               = &mockErr{msg: "EOF"}
	errOther             = &mockErr{msg: "something else"}
)

type mockErr struct{ msg string }

func (e *mockErr) Error() string { return e.msg }
