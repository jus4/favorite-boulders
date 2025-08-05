package handlers

import (
  "context"
  "fmt"
  "github.com/a-h/templ"
  "os"
  "io"
  "net/http"
  "github.com/joho/godotenv"
  "strconv"
  "log"
	"github.com/labstack/echo/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jus4/favorite-boulders/internal/store/queries"
	"github.com/jus4/favorite-boulders/internal/helpers"
	"github.com/jus4/favorite-boulders/internal/templates/components"
)


var noResultsFound = templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
  _, err := w.Write([]byte(""))
  return err
})

var notFoundTxt = templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
  _, err := w.Write([]byte("<p>Ei haku tuloksia</p>"))
  return err
})

func GetRoutesByName(c echo.Context) error {
  keyword := c.FormValue("route-search")
  if len(keyword) <= 2 {
    return helpers.Render(c, http.StatusOK, notFoundTxt)
  }

  routes, err := FetchRoutesByName(keyword)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Haku epäonnistui: %v\n", err)
    return helpers.Render(c, http.StatusInternalServerError, notFoundTxt)
  }


  routeData := make([]components.RouteListProps, len(routes))
  for i, k := range(routes) {
    routeData[i] = components.RouteListProps{
      Name: k.Title.String,
      SectorName: k.SectorName.String,
      SectorId: strconv.FormatInt(k.SectorID, 10),
    }
  }
  return helpers.Render(c, http.StatusOK, components.RouteList(routeData))
}


func SearchRouteByName(c echo.Context) error {
  keyword := c.FormValue("route-search")
  routes, err := FetchRoutesByName(keyword)

  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to get routes: %v\n", err)
    return echo.NewHTTPError(http.StatusInternalServerError, "no results found")
  }


  if len(keyword) <= 2 {
    return helpers.Render(c, http.StatusOK, noResultsFound)
  }

  mappedRoutes := []components.FavouritesListProps{}
  for _, item := range routes {
    mappedRoutes = append(mappedRoutes, components.FavouritesListProps{
      Name: item.Title.String, 
      SectorName: item.SectorName.String, 
      Grade: item.Grade.String,
      ID:strconv.FormatInt(item.RouteID, 10),
    })
  }
  searchResult := components.FavouriteRouteSelector(mappedRoutes)

  return helpers.Render(c, http.StatusOK, searchResult)
}

func SearchEditClimb(c echo.Context) error {
  keyword := c.FormValue("edit-climb-search")
  routes, err := FetchRoutesByName(keyword)

  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to get routes: %v\n", err)
    return echo.NewHTTPError(http.StatusInternalServerError, "no results found")
  }


  if len(keyword) <= 2 {
    return helpers.Render(c, http.StatusOK, noResultsFound)
  }

  mappedRoutes := []components.FavouritesListProps{}
  for _, item := range routes {
    mappedRoutes = append(mappedRoutes, components.FavouritesListProps{
      Name: item.Title.String, 
      SectorName: item.SectorName.String, 
      Grade: item.Grade.String,
      ID:strconv.FormatInt(item.RouteID, 10),
    })
  }
  searchResult := components.SelectEditableRoute(mappedRoutes)

  return helpers.Render(c, http.StatusOK, searchResult)
}

type GetRoutesResponse struct {
  Message string `json:"message"`
  Routes []queries.GetRoutesBySectorRow `json:"routes"`
}

func GetRoutesBySectorId(c echo.Context) error {
  sectorIdStr := c.Param("id")
  sectorName := c.QueryParam("name")
  sectorId, err := strconv.ParseInt(sectorIdStr, 10, 64)
  if err != nil {
    return c.JSON(500, err )
  }
  dbPool := helpers.GetDbPool()
  q := queries.New(dbPool)
  routesData, err := q.GetRoutesBySector(context.Background(), pgtype.Int8{sectorId, true})
  if err != nil {
    return c.JSON(500, err )
  }
  defer dbPool.Close()

  routes := []components.RouteListProps{}
  for _, i := range(routesData) {
    routes = append(routes, components.RouteListProps{
      Name: i.Title.String,
      Grade: i.Grade.String,
      Id: strconv.FormatInt(i.RouteID, 10),
      Image: fmt.Sprint(i.ImageMain),
      RouteType: i.RouteType.String,
    })
  }
  routesRendered := components.SectorRouteList(routes, sectorName)
  if (len(routes) < 1 ) {
    return helpers.Render(c, http.StatusOK, routesRendered)
  }

  return helpers.Render(c, http.StatusOK, routesRendered)

}


func FetchRoutesByName(keyword string) ([]queries.GetRoutesByNameRow, error)  {
  err := godotenv.Load()
  if err != nil {
    log.Print("Error loading .env file")
  }
  conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
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
