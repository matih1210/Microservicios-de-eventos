package di

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "github.com/tuusuario/eventgo/internal/env"
    "github.com/tuusuario/eventgo/internal/domain/event"
)

type Injector struct {
    Cfg *env.Config

    Mongo *mongo.Client

    EventRepo event.Repository
    EventSvc  event.Service
}

func Build(cfg *env.Config) (*Injector, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    cli, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
    if err != nil { return nil, err }

    repo := event.NewMongoRepository(cli, cfg.MongoDB)
    svc := event.NewService(repo)

    return &Injector{
        Cfg:       cfg,
        Mongo:     cli,
        EventRepo: repo,
        EventSvc:  svc,
    }, nil
}

func (i *Injector) Close() {
    if i.Mongo != nil {
        _ = i.Mongo.Disconnect(context.Background())
    }
}
