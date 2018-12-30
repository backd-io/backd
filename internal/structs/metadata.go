package structs

import (
	"time"
)

// Metadata is the struct that represents a metadata information of an struct
type Metadata struct {
	Owner     string `json:"_created_by" bson:"cby"`
	UpdatedBy string `json:"_updated_by" bson:"uby"`
	CreatedAt int64  `json:"_created_at" bson:"cat"`
	UpdatedAt int64  `json:"_updated_at" bson:"uat"`
}

// SetCreate sets the metadata on creation
func (m *Metadata) SetCreate(domainID, userID string) {
	now := time.Now().Unix()
	m.CreatedAt = now
	m.UpdatedAt = now
	m.Owner = FullUsername(domainID, userID)
	m.UpdatedBy = FullUsername(domainID, userID)
}

// SetUpdate sets the metadata on update
func (m *Metadata) SetUpdate(domainID, userID string) {
	m.UpdatedAt = time.Now().Unix()
	m.UpdatedBy = FullUsername(domainID, userID)
}

// FullUsername returns the <domain_id>/<user_id> representation that ensures uniqueness
func FullUsername(domainID, userID string) string {
	return domainID + "/" + userID
}

// MongoDB JSON Schema description for metadata validator
var (
	metadataRequired  = []string{"_meta.cby", "_meta.uby", "_meta.cat", "_meta.uat"}
	metadataValidator = map[string]interface{}{
		"_meta.cby": map[string]interface{}{
			"bsonType": "string",
			"pattern":  "^[a-zA-Z0-9]{20}/[a-zA-Z0-9]{20}$",
		},
		"_meta.uby": map[string]interface{}{
			"bsonType": "string",
			"pattern":  "^[a-zA-Z0-9]{20}/[a-zA-Z0-9]{20}$",
		},
		"_meta.cat": map[string]interface{}{
			"bsonType": "long",
		},
		"_meta.uat": map[string]interface{}{
			"bsonType": "long",
		},
	}
)

// BuildValidator is a helper to join required validations from each collection plus metadata
func BuildValidator(properties map[string]interface{}, required []string) map[string]interface{} {

	for k, v := range metadataValidator {
		properties[k] = v
	}

	for _, item := range metadataRequired {
		required = append(required, item)
	}

	return map[string]interface{}{
		"$jsonSchema": map[string]interface{}{
			"bsonType":   "object",
			"required":   required,
			"properties": properties,
		},
	}

}
