package mongo

import (
    "context"
    "errors"
    "time"

    "github.com/tuusuario/go-auth/internal/domain"
    "github.com/tuusuario/go-auth/internal/repository"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
    client *mongo.Client
}

func NewClient(uri string) (*Client, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        return nil, err
    }
    return &Client{client: c}, nil
}

func (c *Client) Disconnect() {
    _ = c.client.Disconnect(context.Background())
}

type userRepo struct {
    col *mongo.Collection
}

func NewUserRepository(c *Client, dbName string) repository.UserRepository {
    col := c.client.Database(dbName).Collection("users")
    // Unique index on username
    _, _ = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
        Keys: bson.D{{Key: "username", Value: 1}},
        Options: options.Index().SetUnique(true),
    })
    return &userRepo{col: col}
}

func (r *userRepo) Create(u *domain.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _, err := r.col.InsertOne(ctx, u)
    return err
}

func (r *userRepo) FindByUsername(username string) (*domain.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    var u domain.User
    err := r.col.FindOne(ctx, bson.M{"username": username}).Decode(&u)
    if errors.Is(err, mongo.ErrNoDocuments) {
        return nil, nil
    }
    return &u, err
}
