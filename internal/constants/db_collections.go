package constants

// Default Database items naming
const (
	// ColSchema stores information about how the data must be stored and validated
	ColSchema = "_Schemas"
	// ColRoles is the collection that holds the Role Access for every object
	//   managed by the API
	ColRoles = "_Roles"
	// ColRelation is the collection that stores the relations of linked data
	ColRelation = "_Relations"
)

// Domain collections to define security
const (
	ColConfig     = "_Config"     // Contains the especific configuration for this domain
	ColDomains    = "_Domains"    // Contains domains insformation
	ColUsers      = "_Users"      // Contains the users defined at domain level
	ColGroups     = "_Groups"     // Contains the  groups defined at domain level
	ColMembership = "_Membership" // Relationship between users & groups
	ColSession    = "_Sessions"   // Session information (user & group membership), expiration
)

// Default Database for Backd
const (
	DBBackd = "backd"
)
