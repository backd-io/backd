package rest

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/backd-io/backd/internal/instrumentation"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// SetupRouter builds a router for the REST API endpoints
func (rr *REST) SetupRouter(routes map[string]map[string]APIEndpoint, inst *instrumentation.Instrumentation) {

	rr.optionsMappings = make(map[string][]string)

	for method, mappings := range routes {
		for route := range mappings {
			_, ok := rr.optionsMappings[route]
			if ok {
				rr.optionsMappings[route] = append(rr.optionsMappings[route], method)
			} else {
				rr.optionsMappings[route] = []string{method}
			}
		}
	}

	rr.inst = inst
	rr.registerMetrics()

	var router *httprouter.Router
	router = httprouter.New()

	for method, mappings := range routes {
		for route, endpoint := range mappings {

			localMethod := method             // ensure it will be logged
			localRoute := route               // ensure it will be logged
			localFunction := endpoint.Handler // function
			localMatcher := endpoint.Matcher  // ensure naming compliance (if defined)

			// wrapper will handle all logic that are not on the function as:
			//  - instrumentation
			//  - logging
			var wrapper func(w http.ResponseWriter, r *http.Request, ps httprouter.Params)

			wrapper = func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				now := time.Now()
				if match(localRoute, localMatcher, r) {
					ww := NewLogResponseWriter(w)
					rr.writeCORSHeaders(w, r, localRoute)
					localFunction(ww, r, ps)

					rr.log(false, "hit", localMethod, localRoute, r.RequestURI, r.RemoteAddr, ww.Status(), ww.Size(), time.Since(now))
					return
				}

				BadRequest(w, r, "route not match")

				rr.log(true, "hit", localMethod, localRoute, r.RequestURI, r.RemoteAddr, http.StatusBadRequest, 0, time.Since(now))

			}

			router.Handle(method, route, wrapper)

			rr.inst.Info("HTTP route added",
				zap.String("method", localMethod),
				zap.String("route", localRoute),
				zap.String("matchers", strings.Join(localMatcher, ",")),
			)

		}
	}

	// OPTIONS routes
	for route := range rr.optionsMappings {

		localRoute := route // ensure it will be logged
		localMethod := "OPTIONS"
		wrapper := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

			now := time.Now()

			ww := NewLogResponseWriter(w)
			rr.writeCORSHeaders(ww, r, localRoute)
			ww.WriteHeader(http.StatusOK)

			rr.log(false, "hit", localMethod, localRoute, r.RequestURI, r.RemoteAddr, ww.Status(), ww.Size(), time.Since(now))

		}

		router.Handle(localMethod, route, wrapper)

		rr.inst.Info("HTTP route added",
			zap.String("method", localMethod),
			zap.String("route", localRoute),
		)
	}

	// ensure not found and not allowed handlers are logged also
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		ErrorResponse(w, http.StatusNotFound, "")
		rr.log(false, "hit", r.Method, r.RequestURI, r.RequestURI, r.RemoteAddr, http.StatusNotFound, 0, time.Since(now))
	})
	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		ErrorResponse(w, http.StatusMethodNotAllowed, "")
		rr.log(false, "hit", r.Method, r.RequestURI, r.RequestURI, r.RemoteAddr, http.StatusMethodNotAllowed, 0, time.Since(now))
	})

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

	// if there is no matcher to match against then all match
	if len(matcher) == 0 {
		return true
	}

	routeParts = strings.Split(r.URL.Path, "/")
	if len(matcher) != len(routeParts)-1 {
		return false
	}

	for i, m := range matcher {
		// do no check blank or .* routes since everything is already allowed
		if m != "" && m != ".*" {
			matched, err := regexp.MatchString(m, routeParts[i+1])
			if err != nil || matched == false {
				return false
			}
		}
	}

	return true

}

func (rr *REST) writeCORSHeaders(w http.ResponseWriter, r *http.Request, route string) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, X-Session-ID, Connection, Sec-WebSocket-Extensions, Sec-WebSocket-Key, Sec-WebSocket-Version, Upgrade")
	w.Header().Add("Access-Control-Allow-Methods", fmt.Sprintf("OPTIONS, %s", strings.Join(rr.optionsMappings[route], ", ")))
}
