package signup

import (
    "context"
    "errors"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
    Insert(ctx context.Context, s *Signup) (string, error)
    FindActiveByUserAndEvent(ctx context.Context, userID, eventID string) (*Signup, error)
    FindByID(ctx context.Context, id string) (*Signup, error)
    Cancel(ctx context.Context, id string, canceledAt int64) error
    ListActiveByEvent(ctx context.Context, eventID string) ([]Signup, error)
}

type mongoRepo struct {
    col *mongo.Collection
}

func NewMongoRepository(cli *mongo.Client, db string) Repository {
    col := cli.Database(db).Collection("signups")
    _, _ = col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
        { Keys: bson.D{{Key:"eventId",Value:1},{Key:"userId",Value:1},{Key:"canceled",Value:1}} },
        { Keys: bson.D{{Key:"eventId",Value:1},{Key:"canceled",Value:1}} },
    })
    return &mongoRepo{col: col}
}

func (r *mongoRepo) Insert(ctx context.Context, s *Signup) (string, error) {
    s.Created = time.Now().Unix()
    s.Canceled = 0
    res, err := r.col.InsertOne(ctx, s)
    if err != nil { return "", err }
    if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
        return oid.Hex(), nil
    }
    return "", errors.New("invalid inserted id")
}

func (r *mongoRepo) FindActiveByUserAndEvent(ctx context.Context, userID, eventID string) (*Signup, error) {
    var s Signup
    err := r.col.FindOne(ctx, bson.M{"userId": userID, "eventId": eventID, "canceled": 0}).Decode(&s)
    if err == mongo.ErrNoDocuments { return nil, nil }
    return &s, err
}

func (r *mongoRepo) FindByID(ctx context.Context, id string) (*Signup, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil { return nil, err }
    var s Signup
    if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&s); err != nil {
        return nil, err
    }
    s.ID = id
    return &s, nil
}

func (r *mongoRepo) Cancel(ctx context.Context, id string, canceledAt int64) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil { return err }
    _, err = r.col.UpdateByID(ctx, oid, bson.M{"$set": bson.M{"canceled": canceledAt}})
    return err
}

func (r *mongoRepo) ListActiveByEvent(ctx context.Context, eventID string) ([]Signup, error) {
    cur, err := r.col.Find(ctx, bson.M{"eventId": eventID, "canceled": 0}, &options.FindOptions{
        Sort: bson.D{{Key:"created", Value:1}},
        Projection: bson.M{"userName":1,"userId":1,"created":1},
    })
    if err != nil { return nil, err }
    defer cur.Close(ctx)

    var out []Signup
    for cur.Next(ctx) {
        var s Signup
        if err := cur.Decode(&s); err != nil { return nil, err }
        var raw = cur.Current
        if v, err := raw.LookupErr("_id"); err == nil {
            var oid primitive.ObjectID
            if v.Unmarshal(&oid) == nil {
                s.ID = oid.Hex()
            }
        }
        out = append(out, s)
    }
    return out, cur.Err()
}
