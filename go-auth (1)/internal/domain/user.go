package domain

type User struct {
    ID       string `json:"id" bson:"_id,omitempty"`
    Name     string `json:"name" bson:"name"`
    Username string `json:"username" bson:"username"`
    Password string `json:"-" bson:"password"`
    Created  int64  `json:"created" bson:"created"`
}
