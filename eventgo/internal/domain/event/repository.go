package event

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
    Insert(ctx context.Context, e *Event) (string, error)
    Update(ctx context.Context, e *Event) error
    FindByID(ctx context.Context, id string) (*Event, error)
    Cancel(ctx context.Context, id string, canceledAt int64) error
    FindOpen(ctx context.Context, now int64) ([]Event, error)
}

type mongoRepo struct {
    col *mongo.Collection
}

func NewMongoRepository(cli *mongo.Client, db string) Repository {
    col := cli.Database(db).Collection("events")
    _, _ = col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
        { Keys: bson.D{{Key:"when", Value:1}} },
        { Keys: bson.D{{Key:"canceled", Value:1}} },
        { Keys: bson.D{{Key:"ownerId", Value:1}} },
    })
    return &mongoRepo{col: col}
}

func (r *mongoRepo) Insert(ctx context.Context, e *Event) (string, error) {
    e.Created = time.Now().Unix()
    e.Updated = e.Created
    e.Canceled = 0
    res, err := r.col.InsertOne(ctx, e)
    if err != nil { return "", err }
    if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
        return oid.Hex(), nil
    }
    return "", errors.New("invalid inserted id")
}

func (r *mongoRepo) Update(ctx context.Context, e *Event) error {
    oid, err := primitive.ObjectIDFromHex(e.ID)
    if err != nil { return err }
    e.Updated = time.Now().Unix()
    _, err = r.col.UpdateByID(ctx, oid, bson.M{
        "$set": bson.M{ "name": e.Name, "when": e.When, "updated": e.Updated },
    })
    return err
}

func (r *mongoRepo) FindByID(ctx context.Context, id string) (*Event, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil { return nil, err }
    var e Event
    if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&e); err != nil {
        return nil, err
    }
    e.ID = id
    return &e, nil
}

func (r *mongoRepo) Cancel(ctx context.Context, id string, canceledAt int64) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil { return err }
    _, err = r.col.UpdateByID(ctx, oid, bson.M{
        "$set": bson.M{"canceled": canceledAt, "updated": time.Now().Unix()},
    })
    return err
}

func (r *mongoRepo) FindOpen(ctx context.Context, now int64) ([]Event, error) {
    cur, err := r.col.Find(ctx, bson.M{
        "canceled": 0,
        "when": bson.M{"$gt": now},
    }, &options.FindOptions{
        Sort: bson.D{{Key:"when", Value:1}},
        Projection: bson.M{ "name":1, "when":1, "ownerId":1, "ownerName":1 },
    })
    if err != nil { return nil, err }
    defer cur.Close(ctx)

    var out []Event
    for cur.Next(ctx) {
        var e Event
        if err := cur.Decode(&e); err != nil { return nil, err }
        // get _id
        var raw = cur.Current
        if v, err := raw.LookupErr("_id"); err == nil {
            var oid primitive.ObjectID
            if v.Unmarshal(&oid) == nil {
                e.ID = oid.Hex()
            }
        }
        out = append(out, e)
    }
    return out, cur.Err()
}
