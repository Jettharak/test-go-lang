package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadServerConfig(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
		wantErr bool
	}{
		{name: "valid config", content: "server:\n  address: \":9090\"\n", want: ":9090"},
		{name: "empty address", content: "server:\n  address: \"\"\n", wantErr: true},
		{name: "unknown field", content: "server:\n  address: \":9090\"\n  extra: true\n", wantErr: true},
		{name: "malformed YAML", content: "server: [\n", wantErr: true},
		{name: "multiple documents", content: "server:\n  address: \":9090\"\n---\nserver:\n  address: \":9091\"\n", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.yml")
			if err := os.WriteFile(path, []byte(test.content), 0600); err != nil {
				t.Fatal(err)
			}

			config, err := loadServerConfig(path)
			if (err != nil) != test.wantErr {
				t.Fatalf("loadServerConfig() error = %v, wantErr %t", err, test.wantErr)
			}
			if err == nil && config.Server.Address != test.want {
				t.Fatalf("address = %q, want %q", config.Server.Address, test.want)
			}
		})
	}
}

func TestLoadServerConfigMissingFile(t *testing.T) {
	_, err := loadServerConfig(filepath.Join(t.TempDir(), "missing.yml"))
	if err == nil {
		t.Fatal("loadServerConfig() error = nil, want an error")
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{name: "positive numbers", a: 10, b: 20, want: 30},
		{name: "negative numbers", a: -10, b: -20, want: -30},
		{name: "mixed signs", a: -10, b: 20, want: 10},
		{name: "zero", a: 0, b: 0, want: 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := add(test.a, test.b); got != test.want {
				t.Fatalf("add(%d, %d) = %d, want %d", test.a, test.b, got, test.want)
			}
		})
	}
}

func TestAddHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader(`{"num1":10,"num2":20}`))
	response := httptest.NewRecorder()

	addHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got, want := strings.TrimSpace(response.Body.String()), `{"result":30}`; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
}

func TestAddHandlerRejectsInvalidRequests(t *testing.T) {
	tests := []struct {
		name   string
		method string
		body   string
		want   int
		allow  string
	}{
		{name: "wrong method", method: http.MethodGet, body: `{}`, want: http.StatusMethodNotAllowed, allow: http.MethodPost},
		{name: "invalid JSON", method: http.MethodPost, body: `{"num1":`, want: http.StatusBadRequest},
		{name: "non-integer field", method: http.MethodPost, body: `{"num1":"1","num2":2}`, want: http.StatusBadRequest},
		{name: "extra JSON", method: http.MethodPost, body: `{"num1":1,"num2":2} {}`, want: http.StatusBadRequest},
		{name: "malformed trailing JSON", method: http.MethodPost, body: `{"num1":1,"num2":2} {`, want: http.StatusBadRequest},
		{name: "unknown field", method: http.MethodPost, body: `{"num1":1,"num2":2,"extra":3}`, want: http.StatusBadRequest},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/add", strings.NewReader(test.body))
			response := httptest.NewRecorder()

			addHandler(response, request)

			if response.Code != test.want {
				t.Fatalf("status = %d, want %d", response.Code, test.want)
			}
			if test.allow != "" && response.Header().Get("Allow") != test.allow {
				t.Fatalf("Allow = %q, want %q", response.Header().Get("Allow"), test.allow)
			}
		})
	}
}

func TestAddHandlerHandlesResponseWriteError(t *testing.T) {
	writer := &failingResponseWriter{err: errors.New("write failed")}
	request := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader(`{"num1":10,"num2":20}`))

	addHandler(writer, request)

	if !writer.writeCalled {
		t.Fatal("response writer was not called")
	}
}

type failingResponseWriter struct {
	header      http.Header
	err         error
	writeCalled bool
}

func (w *failingResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *failingResponseWriter) Write([]byte) (int, error) {
	w.writeCalled = true
	return 0, w.err
}

func (w *failingResponseWriter) WriteHeader(int) {}
