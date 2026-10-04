package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const modifiedXML = `<?xml version="1.0"?>
<applejuice>
<time time="1"/>
<information xx="1" credits="1234" sessionupload="5" sessiondownload="6" uploadspeed="7" downloadspeed="8" openconnections="9" maxuploadpositions="10" maxuploadpositionsok="1"/>
<shares/>
</applejuice>`

func TestMetricsOutput(t *testing.T) {
	var gotPassword string
	core := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPassword = r.URL.Query().Get("password")
		if r.URL.Path != "/xml/modified.xml" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(modifiedXML))
	}))
	defer core.Close()

	cfg := testConfig(core.URL, "secret")
	recorder := httptest.NewRecorder()
	newHandler(cfg).ServeHTTP(recorder, httptest.NewRequest("GET", "/metrics", nil))

	want := "applejuice_xx 1\napplejuice_credits 1234\napplejuice_sessionupload 5\napplejuice_sessiondownload 6\napplejuice_uploadspeed 7\napplejuice_downloadspeed 8\napplejuice_openconnections 9\napplejuice_maxuploadpositions 10\napplejuice_maxuploadpositionsok 1\n"
	if recorder.Code != 200 || recorder.Body.String() != want {
		t.Fatalf("status %d body %q", recorder.Code, recorder.Body.String())
	}
	if gotPassword != "5ebe2294ecd0e0f08eab7690d2a6ee69" {
		t.Fatalf("password not md5 hashed: %q", gotPassword)
	}
	if recorder.Header().Get("Content-Type") != "text/plain;charset=UTF-8" {
		t.Fatalf("content type %q", recorder.Header().Get("Content-Type"))
	}
}

func TestMetricsErrorHidesPassword(t *testing.T) {
	cfg := testConfig("http://127.0.0.1:1", "secret")
	recorder := httptest.NewRecorder()
	newHandler(cfg).ServeHTTP(recorder, httptest.NewRequest("GET", "/metrics/", nil))

	if recorder.Code != 500 {
		t.Fatalf("status %d", recorder.Code)
	}
	if strings.Contains(recorder.Body.String(), cfg.corePass) {
		t.Fatalf("password leaked: %q", recorder.Body.String())
	}
}

func TestIndex(t *testing.T) {
	recorder := httptest.NewRecorder()
	newHandler(testConfig("http://x", "")).ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))
	if recorder.Code != 200 || !strings.Contains(recorder.Body.String(), "/metrics") {
		t.Fatalf("status %d body %q", recorder.Code, recorder.Body.String())
	}
}

func TestConfig(t *testing.T) {
	env := map[string]string{"CORE_HOST": "192.168.1.2", "CORE_PASSWORD": strings.Repeat("a", 32)}
	cfg := configFromEnv(func(key string) string { return env[key] })
	if cfg.coreHost != "http://192.168.1.2" || cfg.corePort != "9851" || cfg.corePass != env["CORE_PASSWORD"] {
		t.Fatalf("unexpected config %+v", cfg)
	}
	if cfg.timeout != 10*time.Second || cfg.listenAddr != ":80" {
		t.Fatalf("unexpected defaults %+v", cfg)
	}
}

func testConfig(coreURL string, password string) config {
	host := strings.TrimSuffix(coreURL, "/")
	port := "9851"
	if index := strings.LastIndex(host, ":"); index > len("http:") {
		port = host[index+1:]
		host = host[:index]
	}
	return config{coreHost: host, corePort: port, corePass: normalizePassword(password), timeout: 2 * time.Second, listenAddr: ":80"}
}

func TestParseInformationSkipsNamespacedAttributes(t *testing.T) {
	attributes, err := parseInformation([]byte(`<r><information a="1" xmlns:x="urn:x" x:b="2"/></r>`))
	if err != nil || len(attributes) != 1 || attributes[0].Name.Local != "a" {
		t.Fatalf("attributes %+v err %v", attributes, err)
	}
}
