package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func newHandler(cfg config) http.Handler {
	client := &http.Client{Timeout: cfg.timeout}
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, "<a href=/metrics>/metrics</a>")
	})

	metrics := func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")

		body, err := fetch(request.Context(), client, cfg)
		if err != nil {
			log.Printf("core request failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "no response from %s", cfg.safeCoreURL())
			return
		}

		attributes, err := parseInformation(body)
		if err != nil {
			log.Printf("core response invalid: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var output strings.Builder
		for _, attribute := range attributes {
			fmt.Fprintf(&output, "applejuice_%s %s\n", attribute.Name.Local, attribute.Value)
		}
		_, _ = io.WriteString(w, output.String())
	}
	mux.HandleFunc("/metrics", metrics)
	mux.HandleFunc("/metrics/", metrics)
	mux.HandleFunc("/metrics/index.php", metrics)

	return mux
}

func fetch(ctx context.Context, client *http.Client, cfg config) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.coreURL(cfg.corePass), nil)
	if err != nil {
		return nil, redact(err, cfg)
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, redact(err, cfg)
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("core returned status %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, redact(err, cfg)
	}
	return body, nil
}

func redact(err error, cfg config) error {
	message := strings.ReplaceAll(err.Error(), cfg.corePass, "***")
	return fmt.Errorf("%s", message)
}

func parseInformation(body []byte) ([]xml.Attr, error) {
	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = charsetReader

	depth := 0
	rootSeen := false
	informationSeen := false
	var attributes []xml.Attr
	for {
		offset := decoder.InputOffset()
		token, err := decoder.Token()
		if err == io.EOF {
			if !rootSeen {
				return nil, fmt.Errorf("missing XML root")
			}
			return attributes, nil
		}
		if err != nil {
			return nil, fmt.Errorf("parse xml: %w", err)
		}
		switch element := token.(type) {
		case xml.StartElement:
			depth++
			if depth == 1 {
				if rootSeen {
					return nil, fmt.Errorf("multiple XML roots")
				}
				rootSeen = true
			}
			if depth == 2 && !informationSeen && element.Name.Local == "information" && bytes.HasPrefix(body[offset:], []byte("<information")) {
				informationSeen = true
				attributes = plainAttributes(element.Attr)
			}
		case xml.EndElement:
			depth--
		case xml.CharData:
			if depth == 0 && len(bytes.TrimSpace(element)) != 0 {
				return nil, fmt.Errorf("text outside XML root")
			}
		}
	}
}

func charsetReader(label string, input io.Reader) (io.Reader, error) {
	switch strings.ToLower(label) {
	case "utf-8", "utf8", "us-ascii":
		return input, nil
	case "iso-8859-1", "latin1", "latin-1":
		raw, err := io.ReadAll(input)
		if err != nil {
			return nil, err
		}
		runes := make([]rune, len(raw))
		for index, value := range raw {
			runes[index] = rune(value)
		}
		return strings.NewReader(string(runes)), nil
	}
	return nil, fmt.Errorf("unsupported charset %q", label)
}

func plainAttributes(attributes []xml.Attr) []xml.Attr {
	plain := make([]xml.Attr, 0, len(attributes))
	for _, attribute := range attributes {
		if attribute.Name.Space != "" || attribute.Name.Local == "xmlns" {
			continue
		}
		plain = append(plain, attribute)
	}
	return plain
}
