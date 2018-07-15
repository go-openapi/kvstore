package middlewares

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/justinas/alice"
)

// These compression constants are copied from the compress/gzip package.
const (
	encodingGzip = "gzip"

	headerAcceptEncoding  = "Accept-Encoding"
	headerContentEncoding = "Content-Encoding"
	headerContentLength   = "Content-Length"
	headerContentType     = "Content-Type"
	headerVary            = "Vary"
	headerSecWebSocketKey = "Sec-WebSocket-Key"

	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

// gzipResponseWriter is the ResponseWriter that negroni.ResponseWriter is
// wrapped in.
type gzipResponseWriter struct {
	w *gzip.Writer
	http.ResponseWriter
}

// Write writes bytes to the gzip.Writer. It will also set the Content-Type
// header using the net/http library content type detection if the Content-Type
// header was not set yet.
func (grw gzipResponseWriter) Write(b []byte) (int, error) {
	if len(grw.Header().Get(headerContentType)) == 0 {
		grw.Header().Set(headerContentType, http.DetectContentType(b))
	}

	// Delete the content length after we know we have been written to.
	grw.Header().Del(headerContentLength)
	return grw.w.Write(b)
}

// handler struct contains the ServeHTTP method and the compressionLevel to be
// used.
type handler struct {
	compressionLevel int
	next             http.Handler
}

// GzipMW returns a middleware which will handle the Gzip compression in ServeHTTP.
// Valid values for level are identical to those in the compress/gzip package.
func GzipMW(level int) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return Gzip(level, next)
	}
}

// Gzip returns a handler which will handle the Gzip compression in ServeHTTP.
// Valid values for level are identical to those in the compress/gzip package.
func Gzip(level int, next http.Handler) http.Handler {
	return &handler{
		compressionLevel: level,
		next:             next,
	}
}

// ServeHTTP wraps the http.ResponseWriter with a gzip.Writer.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Skip compression if the client doesn't accept gzip encoding.
	if !strings.Contains(r.Header.Get(headerAcceptEncoding), encodingGzip) {
		h.next.ServeHTTP(w, r)
		return
	}

	// Skip compression if client attempt WebSocket connection
	if len(r.Header.Get(headerSecWebSocketKey)) > 0 {
		h.next.ServeHTTP(w, r)
		return
	}

	// Create new gzip Writer. Skip compression if an invalid compression
	// level was set.
	gz, err := gzip.NewWriterLevel(w, h.compressionLevel)
	if err != nil {
		h.next.ServeHTTP(w, r)
		return
	}
	defer gz.Close()

	// Set the appropriate gzip headers.
	headers := w.Header()
	headers.Set(headerContentEncoding, encodingGzip)
	headers.Set(headerVary, headerAcceptEncoding)

	grw := gzipResponseWriter{
		gz,
		w,
	}

	// Call the h.next.ServeHTTP handler supplying the gzipResponseWriter instead of
	// the original.
	h.next.ServeHTTP(grw, r)

}
