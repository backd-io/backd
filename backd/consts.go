package backd

// local client constants
const (
	clientName = "Backd Go Client"
)

// Permission is the required level of permission required to operate
type Permission string

// Headers for the endpoints
const (
	HeaderSessionID     = "X-Session-ID"
	HeaderApplicationID = "X-Application-ID"
)

// Exported permissions
const (
	PermissionRead   Permission = "read"
	PermissionCreate Permission = "create"
	PermissionUpdate Permission = "update"
	PermissionDelete Permission = "delete"
	PermissionAdmin  Permission = "admin"
)

// Session state
const (
	StateAnonymous int = iota
	StateExpired
	StateLoggedIn
)
