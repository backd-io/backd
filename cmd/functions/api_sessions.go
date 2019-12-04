package main

import (
	"context"
	"net/http"
	"time"

	"github.com/backd-io/backd/backd"
	"github.com/backd-io/backd/internal/pbsessions"
	"github.com/backd-io/backd/internal/rest"
)

func (a *apiStruct) getSession(r *http.Request) (session *pbsessions.Session, applicationID string, err error) {

	var (
		cc pbsessions.SessionsClient
	)

	if r.Header.Get(backd.HeaderSessionID) == "" || r.Header.Get(backd.HeaderApplicationID) == "" {
		err = rest.ErrUnauthorized
		return
	}

	cc = pbsessions.NewSessionsClient(a.sessions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err = cc.GetSession(ctx, &pbsessions.GetSessionRequest{
		Id: r.Header.Get(backd.HeaderSessionID),
	})

	applicationID = r.Header.Get(backd.HeaderApplicationID)
	return

}
