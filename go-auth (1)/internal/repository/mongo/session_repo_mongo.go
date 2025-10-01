package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/tuusuario/go-auth/internal/domain"
	"github.com/tuusuario/go-auth/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type sessionRepo struct{ col *mongo.Collection }

func NewSessionRepository(c *Client, db string) repository.SessionRepository {
	return &sessionRepo{col: c.client.Database(db).Collection("sessions")}
}

func (r *sessionRepo) Create(ctx context.Context, s *domain.Session) (string, error) {
	if s == nil || s.ID == "" {
		return "", errors.New("empty session id")
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Guardamos el SID string como _id
	_, err := r.col.InsertOne(ctx, bson.M{
		"_id":     s.ID,
		"userId":  s.UserID,
		"expires": s.Expires,
	})
	if err != nil {
		return "", err
	}
	return s.ID, nil
}

func (r *sessionRepo) FindByID(ctx context.Context, id string) (*domain.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var out domain.Session
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&out)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *sessionRepo) DeleteByID(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *sessionRepo) Exists(ctx context.Context, sid string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	n, err := r.col.CountDocuments(ctx, bson.M{"_id": sid})
	return n > 0, err
}
