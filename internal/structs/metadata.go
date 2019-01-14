package structs

import (
	"time"

	"github.com/backd-io/backd/internal/utils"

	"github.com/mitchellh/mapstructure"
)

// Metadata is the struct that represents a metadata information of an struct
type Metadata struct {
	CreatedBy string `json:"cby" bson:"cby" mapstructure:"cby"`
	UpdatedBy string `json:"uby" bson:"uby" mapstructure:"uby"`
	CreatedAt int64  `json:"cat" bson:"cat" mapstructure:"cat"`
	UpdatedAt int64  `json:"uat" bson:"uat" mapstructure:"uat"`
}

// SetCreate sets the metadata on creation
func (m *Metadata) SetCreate(domainID, userID string) {
	now := time.Now().Unix()
	m.CreatedAt = now
	m.UpdatedAt = now
	m.CreatedBy = FullUsername(domainID, userID)
	m.UpdatedBy = FullUsername(domainID, userID)
}

// SetUpdate sets the metadata on update
func (m *Metadata) SetUpdate(domainID, userID string) {
	m.UpdatedAt = time.Now().Unix()
	m.UpdatedBy = FullUsername(domainID, userID)
}

// FromInterface sets metadata value from a map using mapstructure
func (m *Metadata) FromInterface(meta map[string]interface{}) error {
	return mapstructure.Decode(meta, &m)
}

// FullUsername returns the <domain_id>/<user_id> representation that ensures uniqueness
func FullUsername(domainID, userID string) string {
	return domainID + "/" + userID
}

// MongoDB JSON Schema description for metadata validator
var (
	metadataRequired  = []string{"meta.cby", "meta.uby", "meta.cat", "meta.uat"}
	metadataValidator = map[string]interface{}{
		"meta.cby": map[string]interface{}{
			"bsonType": "string",
			"pattern":  "^[a-zA-Z0-9]+\\/[a-zA-Z0-9]{20}$",
		},
		"meta.uby": map[string]interface{}{
			"bsonType": "string",
			"pattern":  "^[a-zA-Z0-9]+\\/[a-zA-Z0-9]{20}$",
		},
		"meta.cat": map[string]interface{}{
			"bsonType": "long",
		},
		"meta.uat": map[string]interface{}{
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

	utils.Prettify(map[string]interface{}{
		"$jsonSchema": map[string]interface{}{
			"bsonType":   "object",
			"required":   required,
			"properties": properties,
		},
	})

	return map[string]interface{}{
		"$jsonSchema": map[string]interface{}{
			"bsonType":   "object",
			"required":   required,
			"properties": properties,
		},
	}

}

// SchemalessValidator is the validator for a collection without schema validation
func SchemalessValidator() map[string]interface{} {
	return map[string]interface{}{
		"$jsonSchema": map[string]interface{}{
			"bsonType":             "object",
			"additionalProperties": true,
		},
	}
}
