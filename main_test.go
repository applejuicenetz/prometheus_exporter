package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHealthcheckListenAddr(t *testing.T) {
	for _, address := range []string{"127.0.0.1:0", "127.0.0.2:0", "[::1]:0"} {
		t.Run(address, func(t *testing.T) {
			listener, err := net.Listen("tcp", address)
			if err != nil {
				if address == "[::1]:0" {
					t.Skipf("IPv6 unavailable: %v", err)
				}
				t.Fatal(err)
			}
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				if request.Method != http.MethodHead || request.URL.Path != "/" {
					t.Errorf("unexpected request: %s %s", request.Method, request.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
				}
			}))
			_ = server.Listener.Close()
			server.Listener = listener
			server.Start()
			defer server.Close()

			t.Setenv("LISTEN_ADDR", listener.Addr().String())
			t.Setenv("HTTP_PROXY", "http://127.0.0.1:1")
			if code := healthcheck(configFromEnv(os.Getenv).listenAddr); code != 0 {
				t.Fatalf("healthcheck exit %d", code)
			}

			host, port, err := net.SplitHostPort(listener.Addr().String())
			if err != nil {
				t.Fatal(err)
			}
			if host == "127.0.0.1" {
				for _, wildcard := range []string{"", "0.0.0.0"} {
					if code := healthcheck(net.JoinHostPort(wildcard, port)); code != 0 {
						t.Fatalf("wildcard %q exit %d", wildcard, code)
					}
				}
			}
			if host == "::1" {
				if code := healthcheck(net.JoinHostPort("::", port)); code != 0 {
					t.Fatalf("IPv6 wildcard exit %d", code)
				}
			}
		})
	}
}

func TestHealthcheckFailures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	address := server.Listener.Addr().String()
	if code := healthcheck(address); code != 1 {
		t.Fatalf("HTTP 500 exit %d", code)
	}
	server.Close()
	if code := healthcheck(address); code != 1 {
		t.Fatalf("closed server exit %d", code)
	}
	for _, invalid := range []string{"invalid", "localhost:not-a-port"} {
		if code := healthcheck(invalid); code != 1 {
			t.Fatalf("invalid address %q exit %d", invalid, code)
		}
	}
}
