package main

import (
	"bufio"
	"net"
	"net/textproto"
	"net/url"
	"strings"
)

type Request struct {
	Method  string
	Uri     string
	Host    string
	SSL     bool
	Headers map[string]string
	rwc     net.Conn
	brc     *bufio.Reader
}

func parseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	s2 += s1 + 1
	return line[:s1], line[s1+1 : s2], line[s2+1:], true
}

type BadRequestError struct {
	what string
}

func (b *BadRequestError) Error() string {
	return b.what
}

func GetRequest(conn net.Conn) (*Request, error) {
	req := Request{}
	req.rwc = conn
	req.brc = bufio.NewReader(conn)
	tp := textproto.NewReader(req.brc)
	// First line: GET /index.html HTTP/1.0
	requestLine, err := tp.ReadLine()
	if err != nil {
		return &req, err
	}

	method, requestURI, _, ok := parseRequestLine(requestLine)
	if !ok {
		err = &BadRequestError{"malformed HTTP request"}
		return &req, err
	}
	req.Method = method
	req.Uri = requestURI

	// https request
	if method == "CONNECT" {
		req.SSL = true
		req.Uri = "http://" + requestURI
	}

	// get remote host
	uriInfo, err := url.ParseRequestURI(requestURI)
	if err != nil {
		return &req, nil
	}

	// Subsequent lines: Key: value.
	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		return &req, nil
	}

	if uriInfo.Host == "" {
		req.Host = mimeHeader.Get("Host")
	} else {
		if strings.Index(uriInfo.Host, ":") == -1 {
			req.Host = uriInfo.Host + ":80"
		} else {
			req.Host = uriInfo.Host
		}
	}

	headers := make(map[string]string)
	for k, vs := range mimeHeader {
		for _, v := range vs {
			headers[k] = v
		}
	}
	req.Headers = headers
	return &req, nil
}
