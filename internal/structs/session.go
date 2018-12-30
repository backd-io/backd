package structs

// Session is the struct that reflects the information of the user
//   currently logged into the domain
type Session struct {
	ID        string `json:"_id" bson:"_id"`
	Domain    string `json:"domain" bson:"d"`
	User      User   `json:"user" bson:"u"`
	ExpiresAt int64  `json:"expires_at" bson:"eat"`
	Metadata  `json:"_meta" bson:"_meta"`
}

// // TODO: Sessions must not be stored on database
// // Indexes
// var (
// 	SessionIndexes = []Index{
// 		{
// 			Fields: []string{"_id"},
// 			Unique: true,
// 		},
// 	}
// )

// SessionResponse is the struct that will be returned to the client
//   when a session has been established
type SessionResponse struct {
	ID        string `json:"_id"`
	Domain    string `json:"domain"`
	UserID    string `json:"user_id"`
	ExpiresAt int64  `json:"expires_at"`
}
