package event

type Event struct {
    ID         string `json:"id" bson:"_id,omitempty"`
    Name       string `json:"name" bson:"name"`
    When       int64  `json:"when" bson:"when"`
    Updated    int64  `json:"updated" bson:"updated"`
    Created    int64  `json:"created" bson:"created"`
    Canceled   int64  `json:"canceled" bson:"canceled"`
    OwnerID    string `json:"ownerId" bson:"ownerId"`
    OwnerName  string `json:"-" bson:"ownerName"`
}

type EventListItem struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    When      int64  `json:"when"`
    OwnerID   string `json:"ownerId"`
    OwnerName string `json:"ownerName"`
}
