package handlers

import (
  "context"
  "fmt"
  "github.com/a-h/templ"
  "os"
  "io"
  "net/http"
	"github.com/labstack/echo/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jus4/favorite-boulders/internal/store/queries"
	"github.com/jus4/favorite-boulders/internal/helpers"
	"github.com/jus4/favorite-boulders/internal/templates/components"
)

func GetRoutesByName(c echo.Context) error {
  keyword := c.FormValue("route-search")
  routes, err := FetchRoutesByName(keyword)
  notFoundTxt := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
      _, err := w.Write([]byte("<p>No routes found</p>"))
      return err
  })
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to get routes: %v\n", err)
    return helpers.Render(c, http.StatusInternalServerError, notFoundTxt)
  }

  if len(keyword) <= 2 {
    return helpers.Render(c, http.StatusOK, notFoundTxt)
  }

  routeData := make([]components.RouteListProps, len(routes))
  for i, k := range(routes) {
    routeData[i] = components.RouteListProps{
      Name: k.Title.String,
      SectorName: k.SectorName.String,
    }
  }
  return helpers.Render(c, http.StatusOK, components.RouteList(routeData))
}

func FetchRoutesByName(keyword string) ([]queries.GetRoutesByNameRow, error)  {
  dbUrl := "postgres://postgres:hju23e@localhost:5432/favorite_climbs"
  conn, err := pgx.Connect(context.Background(), dbUrl)
  if err != nil {
          fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
          os.Exit(1)
  }
  defer conn.Close(context.Background())

  q := queries.New(conn)
  title := keyword + "%"
  routes, err := q.GetRoutesByName(context.Background(), pgtype.Text{title, true})
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to get routes: %v\n", err)
		return nil, err
  }
  return routes, nil
}
