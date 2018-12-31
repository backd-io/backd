package main

import (
	"net/http"

	"github.com/backd-io/backd/internal/rest"

	"github.com/backd-io/backd/internal/db"
	"github.com/backd-io/backd/internal/instrumentation"
	"github.com/julienschmidt/httprouter"
)

type apiStruct struct {
	inst  *instrumentation.Instrumentation
	mongo *db.Mongo
}

func (a *apiStruct) getSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	rest.Response(w, nil, nil, nil, http.StatusOK, "")
}

func (a *apiStruct) postSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	rest.Response(w, nil, nil, nil, http.StatusOK, "")
}

func (a *apiStruct) deleteSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	rest.Response(w, nil, nil, nil, http.StatusOK, "")
}
