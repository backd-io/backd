package main

import (
	"context"
	"errors"
	"time"

	"github.com/backd-io/backd/cmd/sessions/store"
	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/db"
	"github.com/backd-io/backd/internal/instrumentation"
	"github.com/backd-io/backd/internal/sessionspb"
	"github.com/backd-io/backd/internal/structs"
)

type sessionsServer struct {
	store *store.Store
	inst  *instrumentation.Instrumentation
	mongo *db.Mongo
}

// CreateSession ask for a session creation
func (s sessionsServer) CreateSession(c context.Context, req *sessionspb.CreateSessionRequest) (*sessionspb.Session, error) {

	var (
		session     structs.Session
		user        structs.User
		response    sessionspb.Session
		memberships []map[string]string
		groups      []string
		err         error
	)

	if req.GetDomainId() == "" || req.GetUserId() == "" || req.GetDurationSeconds() == 0 {
		return &response, errors.New("bad")
	}

	err = s.mongo.GetOneByID(req.GetDomainId(), constants.ColUsers, req.GetUserId(), &user)
	if err != nil {
		return &response, err
	}

	err = s.mongo.GetMany(req.GetDomainId(), constants.ColMembership, map[string]interface{}{"u": req.GetUserId()}, []string{}, &memberships, 0, 0)
	if err != nil {
		return &response, err
	}

	for _, membership := range memberships {
		groups = append(groups, membership["g"])
	}

	session.ID = db.NewXID().String()
	session.DomainID = req.DomainId
	session.User = user

	now := time.Now()
	session.CreatedAt = now.Unix()
	session.ExpiresAt = now.Add(time.Duration(req.GetDurationSeconds()) * time.Second).Unix()

	err = s.store.Set(session.ID, session)
	if err != nil {
		return &response, err
	}

	response.Id = session.ID
	response.DomainId = session.DomainID
	response.UserId = user.ID
	response.ExpiresAt = session.ExpiresAt
	response.Groups = groups
	return &response, nil

}

// GetSession returns a session already established
func (s sessionsServer) GetSession(c context.Context, req *sessionspb.GetSessionRequest) (*sessionspb.Session, error) {

	var (
		sess     structs.Session
		response sessionspb.Session
		err      error
	)

	sess, err = s.store.Get(req.GetId())
	if err != nil {
		return &response, err
	}

	if sess.DomainID != req.GetDomainId() || sess.User.ID != req.GetUserId() {
		if err != nil {
			return &response, errors.New("conflict/unauthorized")
		}
	}

	if sess.IsExpired() {
		return &response, errors.New("session expired")
	}

	response.Id = sess.ID
	response.UserId = sess.User.ID
	response.DomainId = sess.DomainID
	response.CreatedAt = sess.CreatedAt
	response.ExpiresAt = sess.ExpiresAt

	return &response, nil

}

// DeleteSession removes a session if exists, returns transaction status as result
func (s sessionsServer) DeleteSession(c context.Context, req *sessionspb.GetSessionRequest) (*sessionspb.Result, error) {

	var (
		sess     structs.Session
		response sessionspb.Result
		err      error
	)

	sess, err = s.store.Get(req.GetId())
	if err != nil {
		return &response, err
	}

	if sess.DomainID != req.GetDomainId() || sess.User.ID != req.GetUserId() {
		if err != nil {
			return &response, errors.New("conflict/unauthorized")
		}
	}

	err = s.store.Delete(req.GetId())
	if err == nil {
		response.Result = true
	}

	return &response, nil

}
