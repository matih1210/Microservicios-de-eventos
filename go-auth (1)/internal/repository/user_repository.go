package repository

import "github.com/tuusuario/go-auth/internal/domain"

type UserRepository interface {
    Create(u *domain.User) error
    FindByUsername(username string) (*domain.User, error)
}
