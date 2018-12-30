package rest

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/backd-io/backd/internal/instrumentation"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// SetupRouter builds a router for the REST API endpoints
func (rr *REST) SetupRouter(routes map[string]map[string]APIHandler, matchers map[string]map[string]APIMatcher, inst *instrumentation.Instrumentation) {

	rr.registerMetrics()

	var router *httprouter.Router
	router = httprouter.New()

	for method, mappings := range routes {
		for route, function := range mappings {

			localFunction := function
			localMethod := method                             // ensure it will be logged
			localRoute := route                               // ensure it will be logged
			localMatcher := matchers[localMethod][localRoute] // ensure meets regexp

			// wrapper will handle all logic that are not on the function as:
			//  - instrumentation
			//  - logging
			var wrapper func(w http.ResponseWriter, r *http.Request, ps httprouter.Params)

			wrapper = func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				now := time.Now()
				if match(localRoute, localMatcher, r) {
					writeCORSHeaders(w, r)
					ww := NewLogResponseWriter(w)
					localFunction(ww, r, ps)

					rr.log(false, "hit", localMethod, r.RequestURI, r.RemoteAddr, ww.Status(), ww.Size(), time.Since(now))
					return
				}
				BadRequest(w, r)

				rr.log(true, "hit", localMethod, r.RequestURI, r.RemoteAddr, 400, 0, time.Since(now))

			}

			router.Handle(method, route, wrapper)

			rr.inst.Info("HTTP route added",
				zap.String("method", localMethod),
				zap.String("route", localRoute),
				zap.String("matchers", strings.Join(localMatcher, ",")),
			)

			// TODO: Add OPTIONS routes
		}
	}

	router.NotFound = http.HandlerFunc(NotFound)

	rr.router = router

}

// Start runs the REST API blocking execution
func (rr *REST) Start() error {
	// graceful stop must be handled by the caller
	rr.httpServer = new(http.Server)
	rr.httpServer.Addr = rr.ipPort
	rr.httpServer.Handler = rr.router
	return rr.httpServer.ListenAndServe()

}

// Shutdown tells the http server to stop gracefully
func (rr *REST) Shutdown() error {
	return rr.httpServer.Shutdown(context.Background())
}

// match is the local function that verifies
func match(route string, matcher []string, r *http.Request) bool {

	var (
		routeParts []string
	)

	routeParts = strings.Split(r.URL.Path, "/")

	if len(matcher) != len(routeParts)-1 {
		return false
	}

	for i, m := range matcher {
		// do no check blank or .* routes since anything is already allowed
		if m != "" && m != ".*" {
			matched, err := regexp.MatchString(m, routeParts[i+1])
			if err != nil || matched == false {
				return false
			}
		}
	}

	return true

}

func writeCORSHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Connection, Sec-WebSocket-Extensions, Sec-WebSocket-Key, Sec-WebSocket-Version, Upgrade")
	w.Header().Add("Access-Control-Allow-Methods", "OPTIONS, HEAD, GET, POST, PUT, DELETE")
}
