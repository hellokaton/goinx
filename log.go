package main

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

func (w *LogResponseWriter) Header() http.Header {
	return w.RW.Header()
}

func (w *LogResponseWriter) Write(data []byte) (s int, err error) {
	s, err = w.RW.Write(data)
	w.Size += s
	return
}

func (w *LogResponseWriter) WriteHeader(r int) {
	w.RW.WriteHeader(r)
	w.RespCode = r
}

func (w *LogResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.RW.(http.Hijacker)
	if ok {
		return hijacker.Hijack()
	}
	return nil, nil, errors.New("webserver doesn't support hijacking")
}
