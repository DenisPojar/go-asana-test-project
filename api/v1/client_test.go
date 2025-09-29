package v1

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

func TestFetchWithRetry(t *testing.T) {
	tests := []struct {
		name     string
		doFunc   func(req *http.Request) (*http.Response, error)
		wantErr  bool
		wantBody string
		maxTries int
	}{
		{
			name: "200 OK returns body",
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(io.Reader(strings.NewReader("success"))),
					Header:     make(http.Header),
				}, nil
			},
			wantErr:  false,
			wantBody: "success",
			maxTries: 3,
		},
		{
			name:     "429 retries then succeeds",
			doFunc:   funcFactory([]int{429, 429, 200}, []string{"", "", "done"}),
			wantErr:  false,
			wantBody: "done",
			maxTries: 5,
		},
		{
			name:     "5xx retries then fails",
			doFunc:   funcFactory([]int{500, 502, 503}, []string{"", "", ""}),
			wantErr:  true,
			wantBody: "",
			maxTries: 3,
		},
		{
			name: "4xx fails immediately",
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 404,
					Body:       io.NopCloser(strings.NewReader("not found")),
					Header:     make(http.Header),
				}, nil
			},
			wantErr:  true,
			wantBody: "",
			maxTries: 3,
		},
		{
			name: "network error retries",
			doFunc: funcFactoryErrors([]error{
				errors.New("network error"),
				errors.New("network error"),
				errors.New("network error"),
			}),
			wantErr:  true, // last attempt returns nil but no body
			wantBody: "",
			maxTries: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &ApiClient{
				client: &mockClient{doFunc: tt.doFunc},
				token:  "token",
			}

			body, err := client.fetchWithRetry("http://example.com", tt.maxTries)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error=%v, got %v", tt.wantErr, err)
			}
			if string(body) != tt.wantBody {
				t.Errorf("expected body=%q, got %q", tt.wantBody, string(body))
			}
		})
	}
}

// Helper to simulate multiple responses with different status codes
func funcFactory(statuses []int, bodies []string) func(req *http.Request) (*http.Response, error) {
	i := 0
	return func(req *http.Request) (*http.Response, error) {
		if i >= len(statuses) {
			return &http.Response{StatusCode: statuses[len(statuses)-1], Body: io.NopCloser(strings.NewReader(bodies[len(bodies)-1])), Header: make(http.Header)}, nil
		}
		resp := &http.Response{
			StatusCode: statuses[i],
			Body:       io.NopCloser(strings.NewReader(bodies[i])),
			Header:     make(http.Header),
		}
		i++
		return resp, nil
	}
}

// Helper to simulate network errors
func funcFactoryErrors(errs []error) func(req *http.Request) (*http.Response, error) {
	i := 0
	return func(req *http.Request) (*http.Response, error) {
		if i >= len(errs) {
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("success")), Header: make(http.Header)}, nil
		}
		err := errs[i]
		i++
		if err != nil {
			return nil, err
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("success")), Header: make(http.Header)}, nil
	}
}
