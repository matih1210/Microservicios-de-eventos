package repository

import (
	"context"

	"github.com/tuusuario/go-auth/internal/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, s *domain.Session) (string, error)
	FindByID(ctx context.Context, id string) (*domain.Session, error)
	DeleteByID(ctx context.Context, id string) error
	Exists(ctx context.Context, sid string) (bool, error)
}
