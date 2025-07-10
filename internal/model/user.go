package model

import (
	"time"
	"github.com/google/uuid"
  "context"
  "github.com/jackc/pgx/v5/pgxpool"
	"github.com/jus4/favorite-boulders/internal/store/queries"
)

type User struct{
  ID        uuid.UUID `json:"id" db:"id"`
  Email     string    `json:"email" db:"email"`
  CreatedAt time.Time `json:"created_at" db:"created_at"`
  UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserRepository interface {
  CreateUser(ctx context.Context, email string)(*User, error) 
  GetUserByEmail(ctx context.Context, email string)(*User, error)
}

type userRepository struct {
    pool *pgxpool.Pool
}

func (r *userRepository) CreateUser(ctx context.Context, email string) (*User, error) {
    q := queries.New(r.pool)
    userData, err := q.CreateUser(ctx, email)
    user := &User{
        ID:        userData.ID.Bytes,
        Email:     userData.Email,
        CreatedAt: userData.CreatedAt.Time,
        UpdatedAt: userData.CreatedAt.Time,
    }

    if err != nil {
        return nil, err
    }
    
    return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {

    q := queries.New(r.pool)
    userData, err := q.GetUserByEmail(ctx, email)
    if err != nil {
        return nil, err
    }

    user := &User{
        ID:        userData.ID.Bytes,
        Email:     userData.Email,
        CreatedAt: userData.CreatedAt.Time,
        UpdatedAt: userData.CreatedAt.Time,
    }

    return user, nil
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
  return &userRepository{pool: pool}
}
