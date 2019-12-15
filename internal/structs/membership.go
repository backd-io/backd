package structs

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Membership is the relation between users and groups
//   While users can have an array of groups or
//     a group can have an array of users it can be degraded if grow
//   Too relation but effective
type Membership struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	UserID  string             `json:"user_id" bson:"user_id"`
	GroupID string             `json:"group_id" bson:"group_id"`
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
			Fields: map[string]interface{}{
				"user_id":  1,
				"group_id": 1,
			},
			Unique: true,
		},
		{
			Fields: map[string]interface{}{"user_id": 1},
			Unique: false,
		},
		{
			Fields: map[string]interface{}{"group_id": 1},
			Unique: false,
		},
	}
)
