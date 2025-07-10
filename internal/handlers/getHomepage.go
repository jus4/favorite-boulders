package handlers

import (
	"net/http"

	"github.com/jus4/favorite-boulders/internal/helpers"
	"github.com/jus4/favorite-boulders/internal/model"
	"github.com/jus4/favorite-boulders/internal/templates/pages"
	"github.com/labstack/echo/v4"
)

func HomePage(c echo.Context) error {
  return helpers.Render(c, http.StatusOK, pages.FrontPage(model.Metadata{Title: "Suomi topot"}, nil))
}
