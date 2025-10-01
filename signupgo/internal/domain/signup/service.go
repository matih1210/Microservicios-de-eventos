package signup

import (
    "context"
    "errors"
    "time"

    "github.com/tuusuario/signupgo/internal/domain/event"
)

var (
    ErrAlreadySigned  = errors.New("already signed up")
    ErrForbidden      = errors.New("forbidden")
    ErrNotFound       = errors.New("not found")
    ErrEventNotFound  = errors.New("event not found")
    ErrEventCanceled  = errors.New("event canceled")
)

type Service interface {
    Create(ctx context.Context, userID, userName, eventID string) (*Signup, error)
    Cancel(ctx context.Context, id, actorID string) error
    ListByEvent(ctx context.Context, eventID string) ([]ListItem, error)
}

type service struct {
    repo   Repository
    events event.Reader
}

func NewService(r Repository, er event.Reader) Service {
    return &service{repo: r, events: er}
}

func (s *service) Create(ctx context.Context, userID, userName, eventID string) (*Signup, error) {
    ev, err := s.events.FindActiveByID(ctx, eventID)
    if err != nil || ev == nil {
        return nil, ErrEventNotFound
    }
    if ev.Canceled != 0 {
        return nil, ErrEventCanceled
    }
    existing, err := s.repo.FindActiveByUserAndEvent(ctx, userID, eventID)
    if err != nil { return nil, err }
    if existing != nil {
        return nil, ErrAlreadySigned
    }
    sgn := &Signup{ UserID: userID, UserName: userName, EventID: eventID }
    id, err := s.repo.Insert(ctx, sgn)
    if err != nil { return nil, err }
    sgn.ID = id
    return sgn, nil
}

func (s *service) Cancel(ctx context.Context, id, actorID string) error {
    sgn, err := s.repo.FindByID(ctx, id)
    if err != nil || sgn == nil {
        return ErrNotFound
    }
    if sgn.UserID != actorID {
        return ErrForbidden
    }
    return s.repo.Cancel(ctx, id, time.Now().Unix())
}

func (s *service) ListByEvent(ctx context.Context, eventID string) ([]ListItem, error) {
    sgns, err := s.repo.ListActiveByEvent(ctx, eventID)
    if err != nil { return nil, err }
    out := make([]ListItem, 0, len(sgns))
    for _, s := range sgns {
        out = append(out, ListItem{
            UserName: s.UserName, UserID: s.UserID, ID: s.ID, SignupDate: s.Created,
        })
    }
    return out, nil
}
