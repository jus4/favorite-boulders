package handlers

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jus4/favorite-boulders/internal/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strconv"
)

func FavouritesCreate(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "could not get session",
		})
	}

	// Access values
	userID, ok := sess.Values["user_email"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	favouritesListName := c.FormValue("favourites_name")
	description := c.FormValue("description")
	selectedClimbs := c.Request().Form["selected_climbs[]"]

	if favouritesListName == "" || len(selectedClimbs) < 1 {
		return c.String(http.StatusOK, "Please fill required fields")
	}

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	selectedRouteIDs := []int32{}
	for _, climb := range selectedClimbs {
		id, err := strconv.Atoi(climb)
		if err != nil {
			return c.String(http.StatusOK, "Invalid route IDs")
		}
		selectedRouteIDs = append(selectedRouteIDs, int32(id))
	}

	userRepository := model.NewFavouritesList(dbpool)
	data, err := userRepository.CreateFavouritesList(context.Background(), struct {
		Name        string
		Description string
		Owner       string
	}{
		Name:        favouritesListName,
		Description: description,
		Owner:       userID,
	})
	// fmt.Print(data.ID)

	if err != nil {
		fmt.Print(err)
		return c.String(http.StatusBadRequest, "There was problem adding the list")
	}

	error := userRepository.InsertFavouriteListItems(context.Background(), data.ID, selectedRouteIDs)
	if error != nil {
		return c.String(http.StatusBadRequest, "There was problem adding routes to the list")
	}

	return c.String(http.StatusOK, "Favourites list added")
}
