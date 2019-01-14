package structs

// Relation is the representation of linked data on the DB.
type Relation struct {
	ID            string `json:"_id" bson:"_id"`
	Source        string `json:"src" bson:"src"`
	SourceID      string `json:"sid" bson:"sid"`
	Destination   string `json:"dst" bson:"dst"`
	DestinationID string `json:"did" bson:"did"`
	Relation      string `json:"rel" bson:"rel"`
	Metadata
}
