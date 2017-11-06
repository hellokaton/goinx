package main

import (
	"log"
	"net/http"
)

type Site struct {
	Name         string    `json:"name"`
	ListenAddr   string    `json:"listen"`
	EnableSSL    bool      `json:"ssl"`
	AddForwarded bool      `json:"add_forwarded"`
	Domains      []*string `json:"domains"`
	ProxyPass    string    `json:"proxy_pass"`
	KeyFile      string    `json:"key_file"`
	CertFile     string    `json:"cert_file"`
}

type Global struct {
	Accesslog string `json:"accesslog"`
}

type Config struct {
	Global Global `json:"global"`
	Sites  []Site `json:"sites"`
}

type RequestHandler struct {
	Transport *http.Transport
	Site      *Site
}

type RequestLogger struct {
	handler http.Handler
	logger  log.Logger
}

type LogResponseWriter struct {
	RW       http.ResponseWriter
	RespCode int
	Size     int
}
