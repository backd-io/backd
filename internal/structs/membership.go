package structs

// Membership is the relation between users and groups
//   While users can have an array of groups or
//     a group can have an array of users it can be degraded if grow
//   Too relation but effective
type Membership struct {
	UserID  string `json:"user_id" bson:"user_id"`
	GroupID string `json:"group_id" bson:"group_id"`
}

// MembershipValidator is a schema for the membership collections
func MembershipValidator() map[string]interface{} {
	return map[string]interface{}{
		"$jsonSchema": map[string]interface{}{
			"bsonType": "object",
			"required": []string{"user_id", "group_id"},
			"properties": map[string]interface{}{
				"user_id": map[string]interface{}{
					"bsonType": "string",
					"pattern":  "^[a-zA-Z0-9]{20}$",
				},
				"group_id": map[string]interface{}{
					"bsonType": "string",
					"pattern":  "^[a-zA-Z0-9]{20}$",
				},
			},
		},
	}
}

// Indexes
var (
	MembershipIndexes = []Index{
		{
			Fields: []string{"user_id", "group_id"},
			Unique: true,
		},
		{
			Fields: []string{"user_id"},
			Unique: false,
		},
		{
			Fields: []string{"group_id"},
			Unique: false,
		},
	}
)
