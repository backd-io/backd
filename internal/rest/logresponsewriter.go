package rest

import "net/http"

// LogResponseWriter is a wrapper to get the information metrics of the
//   http.ResponseWriter
type LogResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// NewLogResponseWriter returns the wrapped http.ResponseWriter
func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w}
}

// Status returns the http error code for the response
func (w *LogResponseWriter) Status() int {
	return w.status
}

// Size returns the amount of data transferred in bytes
func (w *LogResponseWriter) Size() int {
	return w.size
}

// Write wraps the http.ResponseWriter.Write method
func (w *LogResponseWriter) Write(p []byte) (n int, err error) {
	written, err := w.ResponseWriter.Write(p)
	w.size += written

	return written, err
}

// WriteHeader wraps the http.ResponseWriter.WriteHeader method
func (w *LogResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
