package handlers

import (
	"context"
	"fmt"
	"github.com/jus4/favorite-boulders/internal/helpers"
	"github.com/jus4/favorite-boulders/internal/model"
	"github.com/jus4/favorite-boulders/internal/templates/components"
	"github.com/jus4/favorite-boulders/internal/templates/pages"
	"github.com/labstack/echo/v4"
	"net/http"
)

func FavouriteLists(c echo.Context) error {
	userId, ok := c.Get("user_email").(string)
	if !ok {
		return helpers.Render(c, http.StatusOK, pages.FavouriteClimbs(model.Metadata{Title: "Suosikki reitit"}, nil))
	}
	dbpool := helpers.GetDbPool()

	favourites := model.NewFavouritesList(dbpool)
	data, err := favourites.GetFavouriteList(context.Background(), userId)
	if err != nil {
		fmt.Print(err)
	}
	defer dbpool.Close()

	content := components.FavouritesList(data)

	return helpers.Render(c, http.StatusOK, pages.FavouriteClimbs(model.Metadata{Title: "Suosikki reitit"}, content))
}
