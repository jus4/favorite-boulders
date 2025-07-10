package model

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jus4/favorite-boulders/internal/store/queries"
	"golang.org/x/net/context"
	"github.com/jackc/pgx/v5/pgtype"
)

type SectorList struct {
  Latitude    pgtype.Float8
  Longitude   pgtype.Float8
  Name        string
}


type SectorRepository interface {
  GetAllSectors(ctx context.Context)([]SectorList, error)
}

type sectorListRepository struct {
  pool *pgxpool.Pool
}

func ( r *sectorListRepository) GetAllSectors(ctx context.Context)([]SectorList, error) {
  q := queries.New(r.pool)
  sectorData, err := q.GetSectors(ctx)
  if err != nil {
    return nil, err
  }
  var sectors []SectorList
  for _, s := range sectorData {
    sectors = append(sectors, SectorList{Longitude: s.Longitude, Latitude: s.Latitude, Name: s.Name.String}) 
  }
  
  return sectors, nil

}

func NewSectorRepository(pool *pgxpool.Pool) SectorRepository {
  return &sectorListRepository{pool: pool}
}
