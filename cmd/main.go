package main

import (
	"github.com/labstack/echo/v4"
  "github.com/jus4/favorite-boulders/internal/handlers"
  "github.com/labstack/echo-contrib/session"
  "github.com/gorilla/sessions"
  "github.com/jus4/favorite-boulders/internal/middleware"
  "github.com/joho/godotenv"
  "log"
  "net/http"
  "os"
)

func main() {
  err := godotenv.Load()
  if err != nil {
    log.Print("Error loading .env file")
  }
  e := echo.New()
  e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))))
  e.Use(middleware.AuthMiddleware)
  e.GET("/", handlers.HomePage)
  e.GET("/auth/google/login", handlers.OauthGoogleLogin) 
  e.GET("/auth/google/callback", handlers.OauthGoogleCallback) 
  e.GET("/auth/logout", handlers.Logout)
  e.GET("/proxy/wmts/*", func(c echo.Context) error {
	  tilePath := c.Param("*")

	  // Construct full upstream URL
	  targetURL := "https://avoin-karttakuva.maanmittauslaitos.fi/avoin/wmts/" + tilePath + "?api-key=8717d852-4a37-4bd1-a284-83854ffa2478"

	  // Forward the request
	  resp, err := http.Get(targetURL)
	  if err != nil {
	  	return c.String(http.StatusBadGateway, "Failed to fetch tile")
	  }
	  defer resp.Body.Close()

	  return c.Stream(resp.StatusCode, resp.Header.Get("Content-Type"), resp.Body)
  })
  e.POST("/api/get-routes/", handlers.GetRoutesByName)
  e.POST("/api/search-route/", handlers.SearchRouteByName)
  e.POST("/api/favourites-create/", handlers.FavouritesCreate )
  e.GET("/api/get-sectors/", handlers.GetSectors)
  e.GET("/api/routes-by-sector/:id", handlers.GetRoutesBySectorId)
  e.GET("/favourite-climbs/", handlers.FavouriteLists)
  e.Static("/static", "static") 


  // Start the app
  port := os.Getenv("PORT")
  if port == "" {
      port = "8080"
  }

  e.Logger.Fatal(e.Start(":" + port))

}
