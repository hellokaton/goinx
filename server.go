package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	log "github.com/biezhi/goinx/log"
)

type Server struct {
	Name      string   `yaml:"name"`
	Listen    string   `yaml:"listen"`
	Domains   []string `yaml:"domains"`
	Root      *string  `yaml:"root"`
	SSL       bool     `yaml:"ssl"`
	ProxyPass *string  `yaml:"proxy_pass"`
	KeyFile   string   `yaml:"key_file"`
	CertFile  string   `yaml:"cert_file"`
}

func (s *Server) Start() {
	if s.ProxyPass != nil {
		log.Info("%s listen %s, ssl: %v, proxy to %s", s.Name, s.Listen, s.SSL, *s.ProxyPass)
	} else {
		if s.Root != nil {
			log.Info("%s listen %s, ssl: %v, static dir %s", s.Name, s.Listen, s.SSL, *s.Root)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.Handler)

	var err error
	if s.SSL {
		err = http.ListenAndServeTLS(s.Listen, s.CertFile, s.KeyFile, mux)
	} else {
		err = http.ListenAndServe(s.Listen, mux)
	}

	if err != nil {
		log.Error("%v", err)
	}
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", ProjectName+"/"+Version)
	if !Contains(s.Domains, strings.Replace(r.Host, s.Listen, "", -1)) {
		log.Error("Request [%s] RemoteAddr: %s, Header: %v", r.Host, r.RemoteAddr, r.Header)
		w.Write([]byte("Bad Request."))
		return
	}
	log.Info("Request [%s] RemoteAddr: %s, Header: %v", r.Host, r.RemoteAddr, r.Header)

	if s.ProxyPass == nil {
		s.Static(w, r)
	} else {
		s.Proxy(w, r)
	}

}

// static server
func (s *Server) Static(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	log.Info("Request URI: %s", path)
	data, err := ioutil.ReadFile(*s.Root + "/" + string(path))
	if err == nil {
		var contentType string
		if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".html") {
			contentType = "text/html"
		} else if strings.HasSuffix(path, ".js") {
			contentType = "application/javascript"
		} else if strings.HasSuffix(path, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(path, ".svg") {
			contentType = "image/svg+xml"
		} else {
			contentType = "text/plain"
		}

		w.Header().Add("Content Type", contentType)
		w.Write(data)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 My dear - " + http.StatusText(404)))
	}
}

// proxy server
func (s *Server) Proxy(w http.ResponseWriter, r *http.Request) {
	o := new(http.Request)

	*o = *r
	targetURL, err := url.Parse(*s.ProxyPass)

	o.Host = targetURL.Host
	o.URL.Scheme = targetURL.Scheme
	o.URL.Host = targetURL.Host
	o.URL.Path = singleJoiningSlash(targetURL.Path, o.URL.Path)

	if q := o.URL.RawQuery; q != "" {
		o.URL.RawPath = o.URL.Path + "?" + q
	} else {
		o.URL.RawPath = o.URL.Path
	}

	o.URL.RawQuery = targetURL.RawQuery

	o.Proto = "HTTP/1.1"
	o.ProtoMajor = 1
	o.ProtoMinor = 1
	o.Close = false

	transport := http.DefaultTransport

	res, err := transport.RoundTrip(o)

	if err != nil {
		log.Error("http: proxy error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hdr := w.Header()

	for k, vv := range res.Header {
		for _, v := range vv {
			hdr.Add(k, v)
		}
	}

	// for _, c := range res.SetCookie {
	// w.Header().Add("Set-Cookie", c.Raw)
	// }

	w.WriteHeader(res.StatusCode)

	if res.Body != nil {
		io.Copy(w, res.Body)
	}

}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
