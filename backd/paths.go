package backd

import (
	"net/url"
	"path"
)

type microservice int

const (
	admin microservice = iota
	auth
	objects
)

// auth paths
const (
	pathSession string = "session"
)

// admin paths
const (
	pathBootstrap = "bootstrap"
)

func (b *Backd) buildPath(m microservice, route string) string {

	var (
		urlString string
	)

	switch m {
	case admin:
		urlString = b.adminURL
	case auth:
		urlString = b.authURL
	case objects:
		urlString = b.objectsURL
	}

	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, route)
	return u.String()

}
