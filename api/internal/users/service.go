package users

import (
	"context"
	"errors"

	"starterkit/internal/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrUserNotFound = errors.New("user not found")

type Querier interface {
	GetUserByID(ctx context.Context, id pgtype.UUID) (db.GetUserByIDRow, error)
	ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.ListUsersRow, error)
}

type Service struct {
	queries Querier
}

func NewService(queries Querier) *Service {
	return &Service{
		queries: queries,
	}
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	// Convert uuid.UUID to pgtype.UUID
	pgID := pgtype.UUID{}
	if err := pgID.Scan(id.String()); err != nil {
		return nil, err
	}

	dbUser, err := s.queries.GetUserByID(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Convert pgtype.UUID back to uuid.UUID
	var userID uuid.UUID
	if dbUser.ID.Valid {
		userID = uuid.UUID(dbUser.ID.Bytes)
	}

	return &User{
		ID:        userID,
		Email:     dbUser.Email,
		Name:      dbUser.Name,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}, nil
}

func (s *Service) ListUsers(ctx context.Context, limit, offset int) ([]*User, error) {
	// Set default limit if not provided
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	if offset < 0 {
		offset = 0
	}

	dbUsers, err := s.queries.ListUsers(ctx, db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	users := make([]*User, len(dbUsers))
	for i, dbUser := range dbUsers {
		var userID uuid.UUID
		if dbUser.ID.Valid {
			userID = uuid.UUID(dbUser.ID.Bytes)
		}

		users[i] = &User{
			ID:        userID,
			Email:     dbUser.Email,
			Name:      dbUser.Name,
			CreatedAt: dbUser.CreatedAt.Time,
			UpdatedAt: dbUser.UpdatedAt.Time,
		}
	}

	return users, nil
}
