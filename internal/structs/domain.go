package structs

// Domain is a struct that describes the information related to a security domain
//   This information is stored on the `backd` application and its
type Domain struct {
	ID          string `json:"_id" bson:"_id"`
	Name        string `json:"name" bson:"n"`
	Description string `json:"description" bson:"d"`
	Metadata    `json:"_meta" bson:"_meta"`
}

// DomainValidator is the JSON schema validation for the domains collection
func DomainValidator() map[string]interface{} {

	return BuildValidator(
		map[string]interface{}{
			"_id": map[string]interface{}{
				"bsonType": "string",
				"pattern":  "^[a-zA-Z0-9]{20}$",
			},
			"n": map[string]interface{}{
				"bsonType": "string",
			},
			"d": map[string]interface{}{
				"bsonType": "string",
			},
		},
		[]string{"_id", "n"},
	)

}

// Indexes
var (
	DomainIndexes = []Index{
		{
			Fields: []string{"_id"},
			Unique: true,
		},
	}
)
