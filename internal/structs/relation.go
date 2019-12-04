package structs

// Relation is the representation of linked data on the DB.
type Relation struct {
	ID            string `json:"_id" bson:"_id"`
	Source        string `json:"src" bson:"src"`
	SourceID      string `json:"sid" bson:"sid"`
	Destination   string `json:"dst" bson:"dst"`
	DestinationID string `json:"did" bson:"did"`
	Relation      string `json:"rel" bson:"rel"`
	Metadata      `json:"meta" bson:"meta"`
}

// RelationValidator is the JSON schema validation for the domains collection
func RelationValidator() map[string]interface{} {

	return BuildValidator(
		map[string]interface{}{
			"_id": map[string]interface{}{
				"bsonType": "string",
				"pattern":  "^[a-zA-Z0-9]{20}$",
			},
			"src": map[string]interface{}{
				"bsonType": "string",
				"pattern":  "^[a-zA-Z0-9]{2,32}$",
			},
			"sid": map[string]interface{}{
				"bsonType": "string",
				"pattern":  "^[a-zA-Z0-9]{20}$",
			},
			"dst": map[string]interface{}{
				"bsonType": "string",
				"pattern":  "^[a-zA-Z0-9]{2,32}$",
			},
			"did": map[string]interface{}{
				"bsonType": "string",
				"pattern":  "^[a-zA-Z0-9]{20}$",
			},
			"rel": map[string]interface{}{
				"bsonType": "string",
				"pattern":  "^[a-zA-Z0-9]{2,32}$",
			},
		},
		[]string{"_id", "src", "sid", "dst", "did", "rel"},
	)

}

// Indexes
var (
	RelationIndexes = []Index{
		{
			Fields: map[string]interface{}{"_id": 1},
			Unique: true,
		},
		{
			Fields: map[string]interface{}{
				"src": 1,
				"sid": 1,
				"rel": 1,
			},
			Unique: false,
		},
		{
			Fields: map[string]interface{}{
				"dst": 1,
				"did": 1,
				"rel": 1,
			},
			Unique: false,
		},
	}
)
