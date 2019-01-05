package main

import (
	"context"
	"time"

	ldap "gopkg.in/ldap.v2"

	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/pbsessions"
	"github.com/backd-io/backd/internal/structs"
	auth "gopkg.in/korylprince/go-ad-auth.v2"
)

func (a *apiStruct) createSession(sessionRequest structs.SessionRequest) (bool, structs.SessionResponse, error) {

	var (
		sessionResponse structs.SessionResponse
		domain          structs.Domain
		success         bool
		err             error
	)

	// search domain
	err = a.mongo.GetOneByID(constants.DBBackdApp, constants.ColDomains, sessionRequest.DomainID, &domain)
	if err != nil {
		return false, sessionResponse, err
	}

	switch domain.Type {
	case structs.DomainTypeBackd:
		// try to login the user by its password
		var user structs.User
		err = a.mongo.GetOne(domain.ID, constants.ColUsers, map[string]interface{}{"un": sessionRequest.Username}, &user)
		if err != nil {
			return false, sessionResponse, err
		}

		success = user.PasswordMatch(sessionRequest.Password)
		if success {

			var session *pbsessions.Session

			c := pbsessions.NewSessionsClient(a.sessions)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			session, err = c.CreateSession(ctx, &pbsessions.CreateSessionRequest{
				UserId:          user.ID,
				DomainId:        domain.ID,
				DurationSeconds: 300,
				Groups:          []string{},
				External:        false,
			})

			sessionResponse.ExpiresAt = session.GetExpiresAt()
			sessionResponse.ID = session.GetId()

		}

	case structs.DomainTypeActiveDirectory:

		var (
			entry      *ldap.Entry // search entry
			ldapGroups []string    // groups of the user on AD
		)

		success, entry, ldapGroups, err = tryAuthFromAD(sessionRequest.Username, sessionRequest.Password, domain.Config)
		if err != nil {
			return success, sessionResponse, err
		}

		// if the user does not exists create it on the domain DB to set it an ID
		var user structs.User
		err = a.mongo.GetOne(domain.ID, constants.ColUsers, map[string]interface{}{"username": sessionRequest.Username}, &user)
		if err != nil {
			// create user on the DB
			user.Name = entry.GetAttributeValue("displayName")
			user.Email = entry.GetAttributeValue("mail")
			user.Validated = true
			user.Active = true
			err = a.mongo.Insert(domain.ID, constants.ColUsers, &user)
			// if creation fails the return the error
			if err != nil {
				return success, sessionResponse, err
			}
		}

		var session *pbsessions.Session

		c := pbsessions.NewSessionsClient(a.sessions)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		session, err = c.CreateSession(ctx, &pbsessions.CreateSessionRequest{
			UserId:          user.ID,
			DomainId:        domain.ID,
			DurationSeconds: 300,
			Groups:          ldapGroups,
			External:        true,
		})

		sessionResponse.ExpiresAt = session.GetExpiresAt()
		sessionResponse.ID = session.GetId()

	}

	return success, sessionResponse, err
}

func tryAuthFromAD(username, password string, domainConfig map[string]interface{}) (success bool, entry *ldap.Entry, ldapGroups []string, err error) {

	var (
		servers     []string
		port        int
		baseDN      string
		securityInt int
		ok          bool
		config      *auth.Config
	)

	// if we cannot take any of the desired configuration return an error
	if servers, ok = domainConfig["servers"].([]string); !ok {
		err = constants.ErrBadConfiguration
		return
	}

	if port, ok = domainConfig["port"].(int); !ok {
		err = constants.ErrBadConfiguration
		return
	}

	if baseDN, ok = domainConfig["dn"].(string); !ok {
		err = constants.ErrBadConfiguration
		return
	}

	if securityInt, ok = domainConfig["security"].(int); !ok {
		err = constants.ErrBadConfiguration
		return
	}

	for _, server := range servers {
		config = &auth.Config{
			Server:   server,
			Port:     port,
			BaseDN:   baseDN,
			Security: auth.SecurityType(securityInt),
		}

		success, entry, ldapGroups, err = auth.AuthenticateExtended(config, username, password, []string{"displayName", "mail"}, []string{})
		if success {
			return
		}

	}

	return

}
