package backd

// Permission is the required level of permission required to operate
type Permission string

// Exported permissions
const (
	PermissionRead   Permission = "read"
	PermissionCreate Permission = "create"
	PermissionUpdate Permission = "update"
	PermissionDelete Permission = "delete"
	PermissionAdmin  Permission = "admin"
)
