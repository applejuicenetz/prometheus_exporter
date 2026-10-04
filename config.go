package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type config struct {
	coreHost   string
	corePort   string
	corePass   string
	timeout    time.Duration
	listenAddr string
}

func configFromEnv(getenv func(string) string) config {
	timeout := 10 * time.Second
	if seconds, err := strconv.ParseFloat(getenv("PHP_SOCKET_TIMEOUT"), 64); err == nil && seconds > 0 {
		timeout = time.Duration(seconds * float64(time.Second))
	}

	port := getenv("CORE_PORT")
	if port == "" {
		port = "9851"
	}

	listenAddr := getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":80"
	}

	return config{
		coreHost:   normalizeHost(getenv("CORE_HOST")),
		corePort:   port,
		corePass:   normalizePassword(getenv("CORE_PASSWORD")),
		timeout:    timeout,
		listenAddr: listenAddr,
	}
}

func normalizePassword(password string) string {
	if len(password) == 32 {
		return password
	}
	sum := md5.Sum([]byte(password))
	return hex.EncodeToString(sum[:])
}

func normalizeHost(host string) string {
	if strings.HasPrefix(host, "http") {
		return host
	}
	return "http://" + host
}

func (cfg config) coreURL(password string) string {
	return fmt.Sprintf("%s:%s/xml/modified.xml?password=%s", cfg.coreHost, cfg.corePort, password)
}

func (cfg config) safeCoreURL() string {
	return cfg.coreURL("***")
}

func splitPort(listenAddr string) (string, int, error) {
	host, port, err := net.SplitHostPort(listenAddr)
	if err != nil {
		return "", 0, err
	}
	number, err := strconv.Atoi(port)
	return host, number, err
}
