package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = ""
	r.URL.Scheme = "http"

	if h.Site.AddForwarded {
		remote_addr := r.RemoteAddr
		idx := strings.LastIndex(remote_addr, ":")
		if idx != -1 {
			remote_addr = remote_addr[0:idx]
			if remote_addr[0] == '[' && remote_addr[len(remote_addr)-1] == ']' {
				remote_addr = remote_addr[1 : len(remote_addr)-1]
			}
		}
		r.Header.Add("X-Forwarded-For", remote_addr)
	}

	r.URL.Host = h.Site.ProxyPass

	conn_hdr := ""
	conn_hdrs := r.Header["Connection"]

	// log.Printf("Connection headers: %v", conn_hdrs)

	if len(conn_hdrs) > 0 {
		conn_hdr = conn_hdrs[0]
	}

	upgrade_websocket := false
	if strings.ToLower(conn_hdr) == "upgrade" {
		log.Printf("got Connection: Upgrade")
		upgrade_hdrs := r.Header["Upgrade"]
		log.Printf("Upgrade headers: %v", upgrade_hdrs)
		if len(upgrade_hdrs) > 0 {
			upgrade_websocket = (strings.ToLower(upgrade_hdrs[0]) == "websocket")
		}
	}

	if upgrade_websocket {
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
			return
		}

		conn, bufrw, err := hj.Hijack()
		defer conn.Close()

		conn2, err := net.Dial("tcp", r.URL.Host)
		if err != nil {
			http.Error(w, "couldn't connect to backend server", http.StatusServiceUnavailable)
			return
		}
		defer conn2.Close()

		err = r.Write(conn2)
		if err != nil {
			log.Printf("writing WebSocket request to backend server failed: %v", err)
			return
		}
		CopyBidir(conn, bufrw, conn2, bufio.NewReadWriter(bufio.NewReader(conn2), bufio.NewWriter(conn2)))
	} else {
		resp, err := h.Transport.RoundTrip(r)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		for k, v := range resp.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}

		w.WriteHeader(resp.StatusCode)

		io.Copy(w, resp.Body)
		resp.Body.Close()
	}
}

func (s *Site) Start(hosts map[string]string, logger *log.Logger) {
	mux := http.NewServeMux()

	log.Println("Proxy URL:", s.ProxyPass)

	var reqHandler http.Handler = &RequestHandler{
		Transport: &http.Transport{DisableKeepAlives: false, DisableCompression: false},
		Site:      s,
	}

	if logger != nil {
		reqHandler = NewRequestLogger(reqHandler, *logger)
	}

	mux.Handle("/", reqHandler)

	srv := &http.Server{Handler: mux, Addr: s.ListenAddr}

	if s.EnableSSL {
		if err := srv.ListenAndServeTLS(s.CertFile, s.KeyFile); err != nil {
			log.Printf("Starting HTTPS frontend %s failed: %v", s.Name, err)
		}
	} else {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Starting frontend %s failed: %v", s.Name, err)
		}
	}
}

func NewRequestLogger(h http.Handler, l log.Logger) *RequestLogger {
	return &RequestLogger{handler: h, logger: l}
}

func (h *RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request_uri := r.RequestURI
	lrw := &LogResponseWriter{RW: w}
	h.handler.ServeHTTP(lrw, r)
	if lrw.RespCode == 0 {
		lrw.RespCode = 200
	}
	host := "-"
	if r.Host != "" {
		host = r.Host
	}
	h.logger.Printf("%s %s \"%s %s %s\" %d %d", r.RemoteAddr, host, r.Method, request_uri, r.Proto, lrw.RespCode, lrw.Size)
}
