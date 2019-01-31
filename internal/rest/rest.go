package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/backd-io/backd/internal/instrumentation"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// REST holds all logic for the API defined
type REST struct {
	router     *httprouter.Router
	ipPort     string
	inst       *instrumentation.Instrumentation
	httpServer *http.Server
}

// New returns a pointer to a REST struct that holds the interactions for the API
func New(ipPort string) *REST {
	return &REST{
		ipPort: ipPort,
	}
}

func (rr *REST) log(isErr bool, msg, method, route, uri, remote string, httpErrorCode, size int, duration time.Duration) {

	rr.addOperationDuration("backd_rest_ops", method, route, strconv.Itoa(httpErrorCode), duration)
	rr.addOperationCounter("backd_rest_counter", method, route, strconv.Itoa(httpErrorCode))

	if isErr {
		rr.inst.Error(msg,
			zap.String("method", method),
			zap.String("uri", uri),
			zap.String("remote", remote),
			zap.Int("code", httpErrorCode),
			zap.Int("size", size),
			zap.Duration("duration", duration),
		)
		return
	}
	rr.inst.Info(msg,
		zap.String("method", method),
		zap.String("uri", uri),
		zap.String("remote", remote),
		zap.Int("code", httpErrorCode),
		zap.Int("size", size),
		zap.Duration("duration", duration),
	)
}
