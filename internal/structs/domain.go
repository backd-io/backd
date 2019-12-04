package structs

// DomainType defines the behavior to build a session
type DomainType string

const (
	// DomainTypeBackd when set the domain will use natively only the backd users/groups
	DomainTypeBackd DomainType = "b"
	// DomainTypeActiveDirectory when set the domain will inherit the groups from the users
	//   on logon. So user membership will be updated from the ones received when the user
	//   creates a session.
	DomainTypeActiveDirectory DomainType = "ad"
)

// Domain is a struct that describes the information related to a security domain
//   This information is stored on the `backd` application and defines the database
//   that holds the information
type Domain struct {
	ID          string                 `json:"_id" bson:"_id"`
	Description string                 `json:"desc" bson:"desc"`
	Type        DomainType             `json:"type" bson:"type"`
	Config      map[string]interface{} `json:"config,omitempty" bson:"config"`
	Metadata    `json:"meta" bson:"meta"`
}

// DomainValidator is the JSON schema validation for the domains collection
func DomainValidator() map[string]interface{} {

	return BuildValidator(
		map[string]interface{}{
			"_id": map[string]interface{}{
				"bsonType": "string",
				// "pattern":   "^[a-zA-Z0-9]+$",
				"maxLength": 32,
			},
			"desc": map[string]interface{}{
				"bsonType": "string",
			},
			"type": map[string]interface{}{
				"bsonType": "string",
			},
			"c": map[string]interface{}{
				"bsonType": "object",
			},
		},
		[]string{"_id", "type"},
	)

}

// Indexes
var (
	DomainIndexes = []Index{
		{
			Fields: map[string]interface{}{"_id": 1},
			Unique: true,
		},
	}
)
