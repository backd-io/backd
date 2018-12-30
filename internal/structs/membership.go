package structs

// Membership is the relation between users and groups
//   While users can have an array of groups or
//     a group can have an array of users it can be degraded if grow
//   Too relation but effective
type Membership struct {
	UserID  string `json:"user_id" bson:"u"`
	GroupID string `json:"group_id" bson:"g"`
}

// MembershipValidator is a schema for the membership collections
func MembershipValidator() map[string]interface{} {
	return map[string]interface{}{
		"$jsonSchema": map[string]interface{}{
			"bsonType": "object",
			"required": []string{"u", "g"},
			"properties": map[string]interface{}{
				"u": map[string]interface{}{
					"bsonType": "string",
					"pattern":  "^[a-zA-Z0-9]{20}$",
				},
				"g": map[string]interface{}{
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
			Fields: []string{"u", "g"},
			Unique: true,
		},
		{
			Fields: []string{"u"},
			Unique: false,
		},
		{
			Fields: []string{"g"},
			Unique: false,
		},
	}
)
