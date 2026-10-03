package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
	}{
		{name: "wrong method", method: http.MethodGet, body: `{}`, want: http.StatusMethodNotAllowed},
		{name: "invalid JSON", method: http.MethodPost, body: `{"num1":`, want: http.StatusBadRequest},
		{name: "extra JSON", method: http.MethodPost, body: `{"num1":1,"num2":2} {}`, want: http.StatusBadRequest},
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
		})
	}
}
