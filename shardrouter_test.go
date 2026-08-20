package caddyshardrouter

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// nextCalled returns a terminal handler that records whether the middleware
// chain continued past the handler under test.
func nextCalled(called *bool) caddyhttp.HandlerFunc {
	return func(http.ResponseWriter, *http.Request) error {
		*called = true
		return nil
	}
}

// noopNext is a terminal handler for tests that do not assert on whether the
// middleware chain continued.
var noopNext = caddyhttp.HandlerFunc(func(http.ResponseWriter, *http.Request) error { return nil })

// statusOf unwraps the status code a handler signalled, either by returning a
// caddyhttp.HandlerError or by writing directly to the response.
func statusOf(t *testing.T, rec *httptest.ResponseRecorder, err error) int {
	t.Helper()
	if handlerErr, ok := errors.AsType[caddyhttp.HandlerError](err); ok {
		return handlerErr.StatusCode
	}
	if err != nil {
		t.Fatalf("unexpected non-handler error: %v", err)
	}
	return rec.Code
}

func TestBodyShardRouter(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		nilBody    bool
		wantStatus int
		wantBody   string
	}{
		{name: "nil body", nilBody: true, wantStatus: http.StatusNotFound, wantBody: "empty body"},
		{name: "malformed JSON", body: `{"foo": ["bar}`, wantStatus: http.StatusBadRequest, wantBody: "failed to parse JSON"},
		{name: "customer key absent", body: `{"foo": "bar"}`, wantStatus: http.StatusBadRequest, wantBody: "failed to parse customer"},
		{name: "customer not a string", body: `{"customer": 42}`, wantStatus: http.StatusBadRequest, wantBody: "failed to parse customer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			if tt.nilBody {
				r.Body = nil
			}
			rec := httptest.NewRecorder()

			var called bool
			err := BodyShardRouter{}.ServeHTTP(rec, r, nextCalled(&called))

			if got := statusOf(t, rec, err); got != tt.wantStatus {
				t.Errorf("status: got %d, want %d", got, tt.wantStatus)
			}
			if got := rec.Body.String(); !strings.Contains(got, tt.wantBody) {
				t.Errorf("body: got %q, want it to contain %q", got, tt.wantBody)
			}
			if !called {
				t.Error("expected the middleware chain to continue")
			}
		})
	}
}

func TestJWTShardRouter(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		wantStatus int
	}{
		{name: "no Authorization header", authHeader: "", wantStatus: http.StatusUnauthorized},
		{name: "non-bearer scheme", authHeader: "Basic dXNlcjpwYXNz", wantStatus: http.StatusUnauthorized},
		{name: "empty bearer token", authHeader: "Bearer ", wantStatus: http.StatusUnauthorized},
		{name: "wrong signing method", authHeader: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIifQ.UIZchxQD36xuhacrJF9HQ5SIUxH5HBiv9noESAacsxU", wantStatus: http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				r.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			var called bool
			err := JWTShardRouter{}.ServeHTTP(rec, r, nextCalled(&called))

			if got := statusOf(t, rec, err); got != tt.wantStatus {
				t.Errorf("status: got %d, want %d", got, tt.wantStatus)
			}
			if called {
				t.Error("expected the middleware chain to stop")
			}
		})
	}
}

// TestJWTShardRouterOverHTTP exercises the handler through a real client and
// server, over httptest's in-memory network, to confirm the unauthorized
// status survives the round trip.
func TestJWTShardRouterOverHTTP(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := JWTShardRouter{}.ServeHTTP(w, r, noopNext)
		if handlerErr, ok := errors.AsType[caddyhttp.HandlerError](err); ok {
			http.Error(w, handlerErr.Err.Error(), handlerErr.StatusCode)
		}
	})

	server := httptest.NewTestServer(t, handler)

	resp, err := server.Client().Get(server.URL)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body failed: %v", err)
	}
	if want := "missing bearer token"; !strings.Contains(string(body), want) {
		t.Errorf("body: got %q, want it to contain %q", body, want)
	}
}

func TestCaddyModules(t *testing.T) {
	tests := []struct {
		module caddy.Module
		wantID string
	}{
		{module: BodyShardRouter{}, wantID: "http.handlers.body_shard_router"},
		{module: JWTShardRouter{}, wantID: "http.handlers.jwt_shard_router"},
	}

	for _, tt := range tests {
		info := tt.module.CaddyModule()
		if got := string(info.ID); got != tt.wantID {
			t.Errorf("module ID: got %q, want %q", got, tt.wantID)
		}
		if info.New() == nil {
			t.Errorf("%s: New() returned nil", tt.wantID)
		}
	}
}

func TestParseCaddyfile(t *testing.T) {
	if _, err := bodyParseCaddyfile(httpcaddyfile.Helper{}); err != nil {
		t.Errorf("bodyParseCaddyfile: unexpected error: %v", err)
	}
	if _, err := jwtParseCaddyfile(httpcaddyfile.Helper{}); err != nil {
		t.Errorf("jwtParseCaddyfile: unexpected error: %v", err)
	}
}
