package model
import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jus4/favorite-boulders/internal/store/queries"
)


type EditRouteSectorRepository interface {
  GetRouteById(ctx context.Context, id int64)(*queries.GetRouteByIdRow, error)
  UpdateRouteById(ctx context.Context, params queries.UpdateRouteByIdParams)(*queries.UpdateRouteByIdRow, error)
}


type editRouteSectorRepository struct {
    pool *pgxpool.Pool
}

func ( r *editRouteSectorRepository) GetRouteById( ctx context.Context, id int64)(*queries.GetRouteByIdRow, error) {
  q := queries.New(r.pool)
  dbRoute, err := q.GetRouteById(ctx, id)
  if err != nil {
    return nil, err
  }

  return &dbRoute, nil
}

func ( r *editRouteSectorRepository) UpdateRouteById( ctx context.Context, params queries.UpdateRouteByIdParams)(*queries.UpdateRouteByIdRow, error) {
  q := queries.New(r.pool)
  updateRow, updateError := q.UpdateRouteById(ctx, params)
  if updateError != nil {
    return nil, updateError
  }

  return &updateRow, nil
}


func NewEditRouteSectorCrag(pool *pgxpool.Pool) EditRouteSectorRepository {
  return &editRouteSectorRepository{pool: pool}
}
