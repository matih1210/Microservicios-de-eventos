package signup

type Signup struct {
    ID       string `json:"id" bson:"_id,omitempty"`
    UserID   string `json:"userId" bson:"userId"`
    UserName string `json:"userName" bson:"userName"`
    EventID  string `json:"eventId" bson:"eventId"`
    Created  int64  `json:"created" bson:"created"`
    Canceled int64  `json:"canceled" bson:"canceled"`
}

type ListItem struct {
    UserName   string `json:"userName"`
    UserID     string `json:"userId"`
    ID         string `json:"id"`
    SignupDate int64  `json:"signupDate"`
}
