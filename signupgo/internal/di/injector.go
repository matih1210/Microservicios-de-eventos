package di

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tuusuario/signupgo/internal/domain/event"
	"github.com/tuusuario/signupgo/internal/domain/signup"
	"github.com/tuusuario/signupgo/internal/env"
)

type Injector struct {
	Cfg *env.Config

	Mongo *mongo.Client

	SignupRepo signup.Repository
	SignupSvc  signup.Service

	EventReader event.Reader
}

func Build(cfg *env.Config) (*Injector, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, err
	}

	srepo := signup.NewMongoRepository(cli, cfg.MongoDB)
	ereader := event.NewHTTPReader(cfg.EventBaseURL, cfg.HTTPTimeoutMs)

	svc := signup.NewService(srepo, ereader)

	return &Injector{
		Cfg:         cfg,
		Mongo:       cli,
		SignupRepo:  srepo,
		SignupSvc:   svc,
		EventReader: ereader,
	}, nil
}

func (i *Injector) Close() {
	if i.Mongo != nil {
		_ = i.Mongo.Disconnect(context.Background())
	}
}
