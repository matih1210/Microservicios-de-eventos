package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/tuusuario/go-auth/internal/domain"
	"github.com/tuusuario/go-auth/internal/repository"
	"github.com/tuusuario/go-auth/pkg/hash"
)

type UserService interface {
	Register(name, username, password string) (*domain.User, error)
	LoginBasic(username, password string) (string, error)
	CurrentByUsername(username string) (*domain.User, error)
	LogoutBySID(sid string) error
}

type userService struct {
	repoUser  repository.UserRepository
	repoSess  repository.SessionRepository
	jwtSecret string
	jwtExpMin int
}

func NewUserService(
	u repository.UserRepository,
	s repository.SessionRepository,
	secret string,
	expMin int,
) UserService {
	return &userService{
		repoUser:  u,
		repoSess:  s,
		jwtSecret: secret,
		jwtExpMin: expMin,
	}
}

func (s *userService) Register(name, username, password string) (*domain.User, error) {
	// ¿existe username?
	existing, err := s.repoUser.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	// hash de la password y alta
	hashed, err := hash.HashPassword(password)
	if err != nil {
		return nil, err
	}
	u := &domain.User{
		Name:     name,
		Username: username,
		Password: hashed, // guardamos la hash
		Created:  time.Now().Unix(),
	}
	if err := s.repoUser.Create(u); err != nil {
		return nil, err
	}
	u.Password = "" // no exponemos hash
	return u, nil
}

func (s *userService) LoginBasic(username, password string) (string, error) {
	ctx := context.Background()

	u, err := s.repoUser.FindByUsername(username)
	if err != nil || u == nil {
		return "", errors.New("unauthorized")
	}
	if !hash.CheckPassword(u.Password, password) {
		return "", errors.New("unauthorized")
	}

	sid := newSID()
	exp := time.Now().Add(time.Duration(s.jwtExpMin) * time.Minute).Unix()

	claims := jwt.MapClaims{
		"sid": sid,
		"uid": u.ID,
		"usr": u.Username,
		"exp": exp,
		"iat": time.Now().Unix(),
	}
	tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	// 👇 capturamos (string, error) e ignoramos el id devuelto
	if _, err := s.repoSess.Create(ctx, &domain.Session{
		ID:      sid, // guardamos el mismo sid como _id
		UserID:  u.ID,
		Expires: exp,
	}); err != nil {
		return "", err
	}

	return tokenStr, nil
}

// auxiliar: sid aleatorio
func newSID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func (s *userService) CurrentByUsername(username string) (*domain.User, error) {
	u, err := s.repoUser.FindByUsername(username)
	if err != nil || u == nil {
		return nil, errors.New("not found")
	}
	u.Password = ""
	return u, nil
}

func (s *userService) LogoutBySID(sid string) error {
	return s.repoSess.DeleteByID(context.Background(), sid)
}
