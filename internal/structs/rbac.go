package structs

// RBAC is the struct that defines how the role permissions are set on the db
type RBAC struct {
	ID           string `json:"_id" bson:"_id"`
	DomainID     string `json:"domain_id" bson:"did"`
	IdentityID   string `json:"identity_id" bson:"iid"`
	Collection   string `json:"collection" bson:"c"`
	CollectionID string `json:"collection_id" bson:"cid"`
	Permission   string `json:"permission" bson:"p"`
}

// RBACValidator is the JSON schema validation for the domains collection
func RBACValidator() map[string]interface{} {

	return map[string]interface{}{
		"$jsonSchema": map[string]interface{}{
			"bsonType": "object",
			"required": []string{"_id", "did", "iid", "c", "cid", "p"},
			"properties": map[string]interface{}{
				"_id": map[string]interface{}{
					"bsonType": "string",
					"pattern":  "^[a-zA-Z0-9]{20}$",
				},
				"did": map[string]interface{}{
					"bsonType":  "string",
					"pattern":   "[a-zA-Z0-9]+",
					"maxLength": 32,
				},
				"iid": map[string]interface{}{
					"bsonType": "string",
					"pattern":  "^[a-zA-Z0-9]{20}$",
				},
				"c": map[string]interface{}{
					"bsonType":  "string",
					"pattern":   "[a-zA-Z0-9*]+",
					"maxLength": 32,
				},
				"cid": map[string]interface{}{
					"bsonType":  "string",
					"pattern":   "[a-zA-Z0-9*]+",
					"maxLength": 20,
				},
				"p": map[string]interface{}{
					"bsonType":  "string",
					"pattern":   "[a-zA-Z0-9]+",
					"maxLength": 32,
				},
			},
		},
	}

}

// Indexes
var (
	RBACIndexes = []Index{
		{
			Fields: []string{"_id"},
			Unique: true,
		},
		{
			Fields: []string{"did", "iid", "c", "cid", "p"}, // can queries
			Unique: false,
		},
		{
			Fields: []string{"did", "iid", "c", "p"}, // visibleIDs queries
			Unique: false,
		},
	}
)
