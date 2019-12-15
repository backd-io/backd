package structs

// Application is the reference of the application and its configuration.
//  This lives inside the main 'backd' application on _applications collection.
//
type Application struct {
	ID          string `json:"_id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"desc" bson:"desc"`
	Metadata    `json:"meta" bson:"meta"`
}

// ApplicationValidator is the JSON schema validation for the applications collection
func ApplicationValidator() map[string]interface{} {

	return BuildValidator(
		map[string]interface{}{
			"_id": map[string]interface{}{
				"bsonType": "string",
				"pattern":  "^[a-zA-Z0-9]{20}$",
			},
			"name": map[string]interface{}{
				"bsonType": "string",
			},
			"desc": map[string]interface{}{
				"bsonType": "string",
			},
		},
		[]string{"_id", "name"},
	)

}

// Indexes
var (
	ApplicationIndexes = []Index{
		{
			Fields: map[string]interface{}{"_id": 1},
			Unique: true,
		},
	}
)

// var (
// 	ApplicationValidator = map[string]interface{}{
// 		"$jsonSchema": map[string]interface{}{
// 			"bsonType": "object",
// 			"required": []string{"_id", "n"},
// 			"properties": map[string]interface{}{
// 				"_id": map[string]interface{}{
// 					"bsonType": "string",
// 					"pattern":  "^[a-zA-Z0-9]{20}$",
// 				},
// 				"n": map[string]interface{}{
// 					"bsonType": "string",
// 				},
// 				"d": map[string]interface{}{
// 					"bsonType": "string",
// 				},
// 				metadataValidator,
// 			},
// 		},
// 	}
// )
