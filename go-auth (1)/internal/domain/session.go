package domain

type Session struct {
    ID      string `json:"id" bson:"_id,omitempty"`
    UserID  string `json:"userId" bson:"userId"`
    Expires int64  `json:"expires" bson:"expires"`
}
