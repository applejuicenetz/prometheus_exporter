package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const phpReference = `$xml = new SimpleXMLElement(file_get_contents("php://stdin")); foreach ($xml->information->attributes() as $key => $value) { printf("applejuice_%s %s%s", $key, $value, PHP_EOL); }`

func TestPHPCompatibility(t *testing.T) {
	if os.Getenv("PHP_COMPATIBILITY") == "" {
		t.Skip("set PHP_COMPATIBILITY=1")
	}

	samples := map[string]string{
		"core":         `<?xml version="1.0"?><applejuice><time>1</time><information credits="0" sessionupload="1" sessiondownload="2" uploadspeed="3" downloadspeed="4" openconnections="5" maxuploadpositions="50"/></applejuice>`,
		"namespaced":   `<r><information a="1" xmlns:x="urn:x" x:b="2"/></r>`,
		"default ns":   `<r xmlns="urn:x"><information a="1"/></r>`,
		"newline":      "<r><information a=\"x&#10;y\"/></r>",
		"first only":   `<r><information a="1"/><information a="2"/></r>`,
		"nested":       `<r><x><information a="1"/></x></r>`,
		"entities":     `<r><information a="&lt;&amp;&quot;"/></r>`,
		"latin1":       "<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?><r><information a=\"\xe4\"/></r>",
		"no attribute": `<r><information/></r>`,
	}

	for name, sample := range samples {
		t.Run(name, func(t *testing.T) {
			core := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(sample))
			}))
			defer core.Close()

			recorder := httptest.NewRecorder()
			newHandler(testConfig(core.URL, "")).ServeHTTP(recorder, httptest.NewRequest("GET", "/metrics", nil))

			command := exec.Command("docker", "run", "--rm", "-i", "php:8-apache", "php", "-d", "display_errors=0", "-r", phpReference)
			command.Stdin = strings.NewReader(sample)
			want, err := command.Output()
			if err != nil {
				t.Fatalf("PHP reference failed: %v", err)
			}

			if recorder.Code != 200 || recorder.Body.String() != string(want) {
				t.Fatalf("status %d go %q php %q", recorder.Code, recorder.Body.String(), want)
			}
		})
	}
}
