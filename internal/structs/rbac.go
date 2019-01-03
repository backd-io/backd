package structs

// RBAC is the struct that defines how the role permissions are set on the db
type RBAC struct {
	ID           string `json:"_id" bson:"_id"`
	DomainID     string `json:"domain_id" bson:"did"`
	UserID       string `json:"user_id" bson:"uid"`
	Collection   string `json:"collection" bson:"c"`
	CollectionID string `json:"collection_id" bson:"cid"`
	Permission   string `json:"permission" bson:"p"`
}
