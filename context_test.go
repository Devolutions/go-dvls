package dvls

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestContextCancellation(t *testing.T) {
	requestReceived := make(chan struct{})
	allowResponse := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(requestReceived)
		<-allowResponse
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":1,"response":{}}`))
	}))
	defer server.Close()

	client := &Client{
		baseUri: server.URL,
		client:  server.Client(),
		credential: credentials{
			token: "test-token",
		},
	}

	t.Run("RequestWithContext respects cancellation", func(t *testing.T) {
		t.Cleanup(func() { close(allowResponse) })

		ctx, cancel := context.WithCancel(context.Background())

		errCh := make(chan error, 1)
		go func() {
			_, err := client.RequestWithContext(ctx, server.URL+"/test", http.MethodGet, nil)
			errCh <- err
		}()

		<-requestReceived
		cancel()

		err := <-errCh
		if err == nil {
			t.Fatal("expected error due to context cancellation, got nil")
		}
	})
}

func TestContextSucceedsWithoutCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":1,"response":{}}`))
	}))
	defer server.Close()

	client := &Client{
		baseUri: server.URL,
		client:  server.Client(),
		credential: credentials{
			token: "test-token",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.RequestWithContext(ctx, server.URL+"/test", http.MethodGet, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":1,"response":{}}`))
	}))
	defer server.Close()

	client := &Client{
		baseUri: server.URL,
		client:  server.Client(),
		credential: credentials{
			token: "test-token",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.RequestWithContext(ctx, server.URL+"/test", http.MethodGet, nil)
	if err == nil {
		t.Fatal("expected error due to context timeout, got nil")
	}
}

func TestContextPropagation(t *testing.T) {
	type contextKey string
	const testKey contextKey = "test-key"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":1,"response":{"id":"test-id"}}`))
	}))
	defer server.Close()

	contextKeyReceived := false
	originalTransport := http.DefaultTransport
	customTransport := &contextCheckTransport{
		base: originalTransport,
		checkFunc: func(ctx context.Context) {
			if ctx.Value(testKey) == "test-value" {
				contextKeyReceived = true
			}
		},
	}

	client := &Client{
		baseUri: server.URL,
		client:  &http.Client{Transport: customTransport},
		credential: credentials{
			token: "test-token",
		},
	}

	ctx := context.WithValue(context.Background(), testKey, "test-value")
	_, _ = client.RequestWithContext(ctx, server.URL+"/test", http.MethodGet, nil)

	if !contextKeyReceived {
		t.Fatal("context was not properly propagated through the request")
	}
}

type contextCheckTransport struct {
	base      http.RoundTripper
	checkFunc func(context.Context)
}

func (t *contextCheckTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.checkFunc(req.Context())
	return t.base.RoundTrip(req)
}
