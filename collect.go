package gstats

import (
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (r *responseWriter) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.bytesWritten += n
	return n, err
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.statusCode = statusCode
}

type collect struct {
	g *GStats
	h http.Handler
}

func (c *collect) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	w := &responseWriter{
		ResponseWriter: writer,
		statusCode:     200,
		bytesWritten:   0,
	}
	start := time.Now()
	c.h.ServeHTTP(w, request)
	end := time.Now()
	go c.g.notifyRequest(request, w.statusCode, w.bytesWritten, end.Sub(start))
}
