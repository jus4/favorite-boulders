package model

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jus4/favorite-boulders/internal/helpers"
	"github.com/jus4/favorite-boulders/internal/store/queries"
)

type FavoriteList struct {
	ID          pgtype.UUID `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	Description string      `json:"description" db:"description"`
	Owner       string      `json:"owner" db:"owner"`
}

type FavoriteListRepository interface {
	GetFavouriteList(ctx context.Context, id string) ([]FavoriteList, error)
	CreateFavouritesList(ctx context.Context, input struct {
		Name        string
		Description string
		Owner       string
	}) (*FavoriteList, error)
	InsertFavouriteListItems(ctx context.Context, listId pgtype.UUID, routeIds []int32) error
}

type favoriteListRepository struct {
	pool *pgxpool.Pool
}

func (r *favoriteListRepository) GetFavouriteList(ctx context.Context, owner string) ([]FavoriteList, error) {
	q := queries.New(r.pool)
	dbFavourites, err := q.GetFavoriteClimbsList(ctx, owner)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	var favourites []FavoriteList
	for _, f := range dbFavourites {
		favourites = append(favourites, FavoriteList{
			ID:          f.ID,
			Name:        f.Name.String,
			Description: f.Description.String,
			Owner:       f.Owner,
		})
	}

	return favourites, nil
}

func (r *favoriteListRepository) CreateFavouritesList(ctx context.Context, input struct {
	Name        string
	Description string
	Owner       string
}) (*FavoriteList, error) {
	q := queries.New(r.pool)
	favourite, err := q.CreateFavoriteClimbList(ctx, queries.CreateFavoriteClimbListParams{
		Name:        helpers.Text(&input.Name),
		Description: helpers.Text(&input.Description),
		Owner:       input.Owner,
	})

	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	return &FavoriteList{Name: favourite.Name.String, Description: favourite.Description.String, Owner: favourite.Owner, ID: favourite.ID}, nil
}

func (r *favoriteListRepository) InsertFavouriteListItems(ctx context.Context, listId pgtype.UUID, routeIds []int32) error {
	q := queries.New(r.pool)
	_, err := q.InsertFavuriteClimbsItems(ctx, queries.InsertFavuriteClimbsItemsParams{routeIds, listId})

	if err != nil {
		fmt.Print(err)
		return err
	}

	return nil
}

func NewFavouritesList(pool *pgxpool.Pool) FavoriteListRepository {
	return &favoriteListRepository{pool: pool}
}
