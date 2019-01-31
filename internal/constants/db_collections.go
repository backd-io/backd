package constants

// Backd application specific
const (
	ColApplications = "_Applications" // Contains information of the applications managed by the cluster
	ColDomains      = "_Domains"      // Contains domains information and configuration
)

// Application collections
const (
	ColRBAC      = "_RBAC"      // ColRBAC is the collection that holds the Role Access for every object managed by the API
	ColRelations = "_Relations" // ColRelations is the collection that stores the relations of linked data
	ColFunctions = "_Functions" // ColFunctions is the collection that stores the functions code and its configuration
)

// Domain collections to define security
const (
	ColUsers      = "_Users"      // Contains the users defined at domain level
	ColGroups     = "_Groups"     // Contains the groups defined at domain level
	ColMembership = "_Membership" // Relationship between users & groups
)

//  Databases for Backd
const (
	DBBackdApp = "_backd" // this application cannot be browseable by API because of its name. It must not contain any collection or function created by users
	DBBackdDom = "backd"
)
