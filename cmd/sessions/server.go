package main

import (
	"context"
	"errors"
	"time"

	"github.com/backd-io/backd/cmd/sessions/store"
	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/db"
	"github.com/backd-io/backd/internal/instrumentation"
	"github.com/backd-io/backd/internal/pbsessions"
	"github.com/backd-io/backd/internal/structs"
)

type sessionsServer struct {
	store *store.Store
	inst  *instrumentation.Instrumentation
	mongo *db.Mongo
}

// CreateSession ask for a session creation
func (s sessionsServer) CreateSession(c context.Context, req *pbsessions.CreateSessionRequest) (*pbsessions.Session, error) {

	var (
		session     structs.Session
		user        structs.User
		response    pbsessions.Session
		memberships []map[string]string
		groups      []string
		err         error
	)

	if req.GetDomainId() == "" || req.GetUserId() == "" || req.GetDurationSeconds() == 0 {
		return &response, errors.New("bad")
	}

	// if external == true then auth and group membership is managed outside BackD
	//   group membership must be empty
	if req.GetExternal() {
		// search an apply groups if those are defined on the database for role enforcing
		for _, groupName := range req.Groups {
			var group structs.Group
			err = s.mongo.GetOne(req.GetDomainId(), constants.ColGroups, map[string]interface{}{"name": groupName}, &group)
			if err == nil {
				groups = append(groups, group.ID)
			}
		}
	} else {
		// if no group request then expect to have group membership defined
		err = s.mongo.GetOneByID(req.GetDomainId(), constants.ColUsers, req.GetUserId(), &user)
		if err != nil {
			return &response, err
		}

		err = s.mongo.GetMany(req.GetDomainId(), constants.ColMembership, map[string]interface{}{"user_id": req.GetUserId()}, []string{}, &memberships, 0, 0)
		if err != nil {
			return &response, err
		}

		for _, membership := range memberships {
			groups = append(groups, membership["group_id"])
		}

	}

	// build session
	session.ID = db.NewXID().String()
	session.DomainID = req.DomainId
	session.User = user

	now := time.Now()
	session.CreatedAt = now.Unix()
	session.ExpiresAt = now.Add(time.Duration(req.GetDurationSeconds()) * time.Second).Unix()
	session.Groups = groups

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
func (s sessionsServer) GetSession(c context.Context, req *pbsessions.GetSessionRequest) (*pbsessions.Session, error) {

	var (
		sess     structs.Session
		response pbsessions.Session
		err      error
	)

	sess, err = s.store.Get(req.GetId())
	if err != nil {
		return &response, err
	}

	if sess.IsExpired() {
		return &response, errors.New("session expired")
	}

	response.Id = sess.ID
	response.UserId = sess.User.ID
	response.DomainId = sess.DomainID
	response.CreatedAt = sess.CreatedAt
	response.ExpiresAt = sess.ExpiresAt
	response.Groups = sess.Groups
	return &response, nil

}

// DeleteSession removes a session if exists, returns transaction status as result
func (s sessionsServer) DeleteSession(c context.Context, req *pbsessions.GetSessionRequest) (*pbsessions.Result, error) {

	var (
		response pbsessions.Result
		err      error
	)

	_, err = s.store.Get(req.GetId())
	if err != nil {
		return &response, err
	}

	err = s.store.Delete(req.GetId())
	if err == nil {
		response.Result = true
	}

	return &response, nil

}
