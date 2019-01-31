package main

import (
	"fmt"
	"net/http"

	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/pbsessions"
	"github.com/backd-io/backd/internal/rest"
	"github.com/backd-io/backd/internal/structs"
	"github.com/backd-io/backd/pkg/lua"
	"github.com/julienschmidt/httprouter"
)

// POST /functions/:name
func (a *apiStruct) postFunction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		input         map[string]interface{}
		output        map[string]interface{}
		session       *pbsessions.Session
		function      structs.Function
		applicationID string
		luaInstance   *lua.Lang
		err           error
	)

	// getSession & rbac
	session, applicationID, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	err = rest.GetFromBody(r, &input)
	fmt.Println("GetFromBody.err:", err)
	fmt.Println("GetFromBody.input:", input)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	err = a.mongo.GetOne(applicationID, constants.ColFunctions, map[string]string{
		"name": ps.ByName("name"),
	}, &function)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	luaInstance = a.lua.Clone()
	luaInstance.SetAppID(applicationID)
	luaInstance.SetSession(session.GetId(), session.GetExpiresAt())

	output, err = luaInstance.RunFunction(function.Code, input)

	rest.Response(w, output, err, 200, "")

}
