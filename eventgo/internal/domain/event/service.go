package event

import (
    "context"
    "errors"
    "time"
)

var (
    ErrNotFound = errors.New("not found")
    ErrForbidden = errors.New("forbidden")
    ErrPastDate = errors.New("when cannot be in the past")
)

type Service interface {
    Create(ctx context.Context, name string, when int64, ownerID, ownerName string) (*Event, error)
    Update(ctx context.Context, id, name string, when int64, actorID string) (*Event, error)
    Cancel(ctx context.Context, id string, actorID string) error
    GetOpen(ctx context.Context) ([]EventListItem, error)
    GetByID(ctx context.Context, id string) (*Event, error)
}

type service struct {
    repo Repository
}

func NewService(r Repository) Service { return &service{repo: r} }

func (s *service) Create(ctx context.Context, name string, when int64, ownerID, ownerName string) (*Event, error) {
    now := time.Now().Unix()
    if when <= now {
        return nil, ErrPastDate
    }
    e := &Event{
        Name: name,
        When: when,
        OwnerID: ownerID,
        OwnerName: ownerName,
    }
    id, err := s.repo.Insert(ctx, e)
    if err != nil { return nil, err }
    e.ID = id
    return e, nil
}

func (s *service) Update(ctx context.Context, id, name string, when int64, actorID string) (*Event, error) {
    now := time.Now().Unix()
    if when <= now {
        return nil, ErrPastDate
    }
    e, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, ErrNotFound
    }
    if e.OwnerID != actorID {
        return nil, ErrForbidden
    }
    e.Name = name
    e.When = when
    if err := s.repo.Update(ctx, e); err != nil { return nil, err }
    e.Updated = time.Now().Unix()
    return e, nil
}

func (s *service) Cancel(ctx context.Context, id string, actorID string) error {
    e, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return ErrNotFound
    }
    if e.OwnerID != actorID {
        return ErrForbidden
    }
    return s.repo.Cancel(ctx, id, time.Now().Unix())
}

func (s *service) GetOpen(ctx context.Context) ([]EventListItem, error) {
    now := time.Now().Unix()
    es, err := s.repo.FindOpen(ctx, now)
    if err != nil { return nil, err }
    out := make([]EventListItem, 0, len(es))
    for _, e := range es {
        out = append(out, EventListItem{
            ID: e.ID, Name: e.Name, When: e.When,
            OwnerID: e.OwnerID, OwnerName: e.OwnerName,
        })
    }
    return out, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*Event, error) {
    e, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, ErrNotFound
    }
    return e, nil
}
